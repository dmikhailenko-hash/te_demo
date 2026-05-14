package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const DefaultBaseURL = "https://test-clientapi.traderevolution.com/traderevolution/v1"

type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

func NewClient(baseURL, token string) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		token:   token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetToken fetches token via POST /authorize?login=&password=
func GetToken(baseURL, username, password string) (string, error) {
	authBase := toWebhooksHost(strings.TrimRight(baseURL, "/"))
	u := fmt.Sprintf("%s/authorize?login=%s&password=%s",
		authBase,
		url.QueryEscape(username),
		url.QueryEscape(password),
	)

	client := &http.Client{Timeout: 15 * time.Second}
	req, _ := http.NewRequest("POST", u, nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, truncate(string(body), 200))
	}

	var result struct {
		S string `json:"s"`
		D struct {
			AccessToken string `json:"access_token"`
		} `json:"d"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("parse: %v (body: %s)", err, truncate(string(body), 200))
	}
	if result.D.AccessToken != "" {
		return result.D.AccessToken, nil
	}
	return "", fmt.Errorf("no access_token in response: %s", truncate(string(body), 200))
}

func toWebhooksHost(base string) string {
	if strings.Contains(base, "webhooks-clientapi") {
		return base
	}
	if strings.Contains(base, "test-clientapi") {
		return strings.Replace(base, "test-clientapi", "webhooks-clientapi", 1)
	}
	return base
}

func (c *Client) Get(path string, params url.Values) (json.RawMessage, error) {
	u := c.baseURL + path
	if len(params) > 0 {
		u += "?" + params.Encode()
	}
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	return c.do(req)
}

func (c *Client) Post(path string, body interface{}) (json.RawMessage, error) {
	var r io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		r = bytes.NewReader(b)
	}
	req, err := http.NewRequest("POST", c.baseURL+path, r)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return c.do(req)
}

func (c *Client) Patch(path string, body interface{}) (json.RawMessage, error) {
	b, _ := json.Marshal(body)
	req, err := http.NewRequest("PATCH", c.baseURL+path, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return c.do(req)
}

func (c *Client) Delete(path string, body interface{}) (json.RawMessage, error) {
	var r io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		r = bytes.NewReader(b)
	}
	req, err := http.NewRequest("DELETE", c.baseURL+path, r)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return c.do(req)
}

func (c *Client) do(req *http.Request) (json.RawMessage, error) {
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request: %v", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == 401 {
		return nil, fmt.Errorf("unauthorized (401) — set new token: te_demo config --token <token>")
	}
	if resp.StatusCode >= 400 {
		var errResp struct {
			ErrMsg string `json:"errmsg"`
		}
		if json.Unmarshal(body, &errResp) == nil && errResp.ErrMsg != "" {
			return nil, fmt.Errorf("%s", errResp.ErrMsg)
		}
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, truncate(string(body), 300))
	}
	return json.RawMessage(body), nil
}

func truncate(s string, n int) string {
	if len(s) > n {
		return s[:n] + "..."
	}
	return s
}
