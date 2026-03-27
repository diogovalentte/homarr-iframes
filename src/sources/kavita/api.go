package kavita

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (k *Kavita) baseRequest(method, url string, body io.Reader, target any) error {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", k.Token))

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	resBody, err := io.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusUnauthorized {
		err := k.RefreshCurrentToken()
		if err != nil {
			return fmt.Errorf("error refreshing token: %w", err)
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", k.Token))
		resp, err = client.Do(req)
		if err != nil {
			return fmt.Errorf("request status (%s): %s", resp.Status, string(resBody))
		}
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request status (%s): %s", resp.Status, string(resBody))
	}

	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
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

	k.Token = loginResponse.Token
	k.RefreshToken = loginResponse.RefreshToken

	return nil
}

func (k *Kavita) RefreshCurrentToken() error {
	if k.RefreshToken == "" || k.Token == "" {
		return fmt.Errorf("refresh token or token not set")
	}
	loginBody := map[string]string{
		"token":        k.Token,
		"refreshToken": k.RefreshToken,
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

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error: %s", resp.Status)
	}

	defer resp.Body.Close()

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	var loginResponse LoginResponse
	if err := json.Unmarshal(resBody, &loginResponse); err != nil {
		return fmt.Errorf("error unmarshaling JSON: %s\n. Reponse body: %s", err.Error(), string(resBody))
	}

	k.Token = loginResponse.Token
	k.RefreshToken = loginResponse.RefreshToken

	return nil
}

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}
