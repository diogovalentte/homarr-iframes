package pihole

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func (p *Pihole) baseRequest(method, url string, body io.Reader, target any, unauthorizedRetries int) error {
	// buffer the body so the retry below can replay it
	var reqBody []byte
	if body != nil {
		var err error
		reqBody, err = io.ReadAll(body)
		if err != nil {
			return fmt.Errorf("error reading request body: %w", err)
		}
	}

	return p.doRequest(method, url, reqBody, target, unauthorizedRetries)
}

func (p *Pihole) doRequest(method, url string, reqBody []byte, target any, unauthorizedRetries int) error {
	client := &http.Client{}
	var payload io.Reader
	if reqBody != nil {
		payload = bytes.NewReader(reqBody)
	}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// the session used by the request, so the retry below only renews the
	// session it actually failed with. Empty when authenticating with a <v6.0
	// token, which has no session at all.
	var usedSID string
	if p.Password != "" {
		sid, validityTime := p.getSession()
		if sid == "" || time.Now().After(validityTime) {
			if err := p.renewSession(sid); err != nil {
				return err
			}
			sid, _ = p.getSession()
		}
		usedSID = sid
		req.Header.Set("X-FTL-SID", sid)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var jsonErr struct {
			Error struct {
				Message string `json:"message"`
			} `json:"error"`
		}
		if err := json.Unmarshal(resBody, &jsonErr); err == nil {
			if jsonErr.Error.Message == "Unauthorized" && unauthorizedRetries > 0 {
				if err := p.renewSession(usedSID); err != nil {
					return err
				}
				return p.doRequest(method, url, reqBody, target, unauthorizedRetries-1)
			}
		}
		return fmt.Errorf("request status (%s): %s", resp.Status, string(resBody))
	}

	if err := json.Unmarshal(resBody, target); err != nil {
		return fmt.Errorf("error unmarshaling JSON: %s\n. Reponse body: %s", err.Error(), string(resBody))
	}

	return nil
}

func (p *Pihole) Login() error {
	loginBody := map[string]string{
		"password": p.Password,
	}
	jsonData, err := json.Marshal(loginBody)
	if err != nil {
		return err
	}
	payload := bytes.NewReader(jsonData)

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/auth", p.InternalAddress), payload)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("(%s) error reading response body: %w", resp.Status, err)
	}

	if resp.StatusCode != http.StatusOK {
		if strings.Contains(string(resBody), "No 2FA token found in JSON payload") {
			return fmt.Errorf("authentication failed: 2FA is enabled. Use app password instead of your account password")
		}
		return fmt.Errorf("error: %s / Response body: %s", resp.Status, string(resBody))
	}

	var loginResponse LoginResponse
	if err := json.Unmarshal(resBody, &loginResponse); err != nil {
		return fmt.Errorf("error unmarshaling JSON: %s\n. Reponse body: %s", err.Error(), string(resBody))
	}

	if !loginResponse.Session.Valid {
		return fmt.Errorf("authentication failed: invalid session")
	}

	p.setSession(loginResponse.Session.SID, time.Now().Add(time.Duration(loginResponse.Session.Validity-5)*time.Second))

	return nil
}

func (p *Pihole) Logout() error {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/auth", p.InternalAddress), nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	sid, _ := p.getSession()
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-FTL-SID", sid)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("(%s) error reading response body: %w", resp.Status, err)
	}

	if resp.StatusCode != http.StatusNoContent {
		var jsonErr struct {
			Session struct {
				Message string `json:"message"`
			} `json:"session"`
		}
		if err := json.Unmarshal(resBody, &jsonErr); err == nil {
			if jsonErr.Session.Message != "session unknown" {
				return fmt.Errorf("error: %s / Response body: %s", resp.Status, string(resBody))
			}
		} else {
			return fmt.Errorf("error: %s / Response body: %s", resp.Status, string(resBody))
		}
	}

	p.setSession("", time.Time{})

	return nil
}

type LoginResponse struct {
	Session struct {
		Valid    bool   `json:"valid"`
		TOTP     bool   `json:"totp"`
		SID      string `json:"sid"`
		CSRF     string `json:"csrf"`
		Validity int    `json:"validity"`
	} `json:"session"`
}
