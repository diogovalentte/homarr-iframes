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

func (p *Pihole) baseRequest(method, url string, body io.Reader, target interface{}) error {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if p.SID != "" {
		if time.Now().After(p.ValidityTime) {
			err = p.Logout()
			if err != nil {
				return err
			}
			err = p.Login()
			if err != nil {
				return err
			}
		}
		req.Header.Set("X-FTL-SID", p.SID)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	resBody, err := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error: %s | %s", resp.Status, string(resBody))
	}

	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
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

	p.SID = loginResponse.Session.SID
	p.ValidityTime = time.Now().Add(time.Duration(loginResponse.Session.Validity-5) * time.Second)

	return nil
}

func (p *Pihole) Logout() error {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/auth", p.InternalAddress), nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-FTL-SID", p.SID)

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

	p.SID = ""
	p.ValidityTime = time.Time{}

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
