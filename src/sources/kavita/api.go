package kavita

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (k *Kavita) baseRequest(method, url string, body io.Reader, target any) error {
	// buffer the body so the request can be replayed after refreshing the token
	var reqBody []byte
	if body != nil {
		var err error
		reqBody, err = io.ReadAll(body)
		if err != nil {
			return fmt.Errorf("error reading request body: %w", err)
		}
	}
	// the token used by the request, so the retry below only refreshes the token
	// it actually failed with
	var usedToken string
	newRequest := func() (*http.Request, error) {
		var payload io.Reader
		if reqBody != nil {
			payload = bytes.NewReader(reqBody)
		}
		req, err := http.NewRequest(method, url, payload)
		if err != nil {
			return nil, fmt.Errorf("error creating request: %w", err)
		}
		usedToken, _ = k.getTokens()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", usedToken))

		return req, nil
	}

	client := &http.Client{}
	req, err := newRequest()
	if err != nil {
		return err
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

	if resp.StatusCode == http.StatusUnauthorized {
		if err := k.refreshToken(usedToken); err != nil {
			return fmt.Errorf("error refreshing token: %w", err)
		}

		req, err = newRequest()
		if err != nil {
			return err
		}

		resp, err = client.Do(req)
		if err != nil {
			return fmt.Errorf("error sending request: %w", err)
		}
		defer resp.Body.Close()

		resBody, err = io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error reading response body: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("request status (%s): %s", resp.Status, string(resBody))
		}
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request status (%s): %s", resp.Status, string(resBody))
	}

	if err := json.Unmarshal(resBody, target); err != nil {
		return fmt.Errorf("error unmarshaling JSON: %s\n. Reponse body: %s", err.Error(), string(resBody))
	}

	return nil
}

func (k *Kavita) Login() error {
	loginBody := map[string]string{
		"username": k.Username,
		"password": k.Password,
	}
	jsonData, err := json.Marshal(loginBody)
	if err != nil {
		return err
	}
	payload := bytes.NewReader(jsonData)

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/account/login", k.InternalAddress), payload)
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
		return fmt.Errorf("error: %s / Response body: %s", resp.Status, string(resBody))
	}

	var loginResponse LoginResponse
	if err := json.Unmarshal(resBody, &loginResponse); err != nil {
		return fmt.Errorf("error unmarshaling JSON: %s\n. Reponse body: %s", err.Error(), string(resBody))
	}

	if loginResponse.Token == "" || loginResponse.RefreshToken == "" {
		return fmt.Errorf("login failed: token or refresh token is empty, response body: %s", string(resBody))
	}

	k.setTokens(loginResponse.Token, loginResponse.RefreshToken)

	return nil
}

// RefreshCurrentToken refreshes the current token, logging in again if the
// refresh token is gone or no longer accepted.
func (k *Kavita) RefreshCurrentToken() error {
	return k.refreshToken("")
}

// refreshToken refreshes the current token. If staleToken isn't empty and the
// current token is already a different one, another request refreshed it in the
// meantime and there is nothing to do.
func (k *Kavita) refreshToken(staleToken string) error {
	k.refreshMutex.Lock()
	defer k.refreshMutex.Unlock()

	token, refreshToken := k.getTokens()
	if staleToken != "" && token != staleToken {
		return nil
	}

	if refreshToken == "" || token == "" {
		// nothing to refresh with, fall back to a full login so the client can
		// recover on its own instead of staying broken until a restart
		return k.Login()
	}
	loginBody := map[string]string{
		"token":        token,
		"refreshToken": refreshToken,
	}
	jsonData, err := json.Marshal(loginBody)
	if err != nil {
		return err
	}
	payload := bytes.NewReader(jsonData)

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/account/refresh-token", k.InternalAddress), payload)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		// the refresh token expired, a full login is the only way out
		if err := k.Login(); err != nil {
			return fmt.Errorf("error logging in: %w", err)
		}

		return nil
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error: %s", resp.Status)
	}

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	var loginResponse LoginResponse
	if err := json.Unmarshal(resBody, &loginResponse); err != nil {
		return fmt.Errorf("error unmarshaling JSON: %s\n. Reponse body: %s", err.Error(), string(resBody))
	}

	if loginResponse.Token == "" || loginResponse.RefreshToken == "" {
		// never overwrite the current tokens with empty ones, it'd leave the
		// client permanently unable to authenticate. Log in again instead.
		if err := k.Login(); err != nil {
			return fmt.Errorf("refresh response doesn't have the tokens (body: %s), error logging in: %w", string(resBody), err)
		}

		return nil
	}

	k.setTokens(loginResponse.Token, loginResponse.RefreshToken)

	return nil
}

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}
