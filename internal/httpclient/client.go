package httpclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"
)

type AuthTransport struct {
	Username string
	Password string
	Base     http.RoundTripper
}

var _ http.RoundTripper = (*AuthTransport)(nil)

func (t *AuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(t.Username, t.Password)
	req.Header.Set("User-Agent", "curl/7.68.0")
	req.Header.Set("Accept", "application/json")
	return t.Base.RoundTrip(req)
}

type CustomHTTPClient struct {
	BaseURL string
	HTTP    *http.Client
}

func NewCustomHTTPClient(baseUrl, username, password string) *CustomHTTPClient {
	return &CustomHTTPClient{
		BaseURL: baseUrl,
		HTTP: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &AuthTransport{
				Username: username,
				Password: password,
				Base:     http.DefaultTransport,
			},
		},
	}
}

func (c *CustomHTTPClient) MakeGETRequest(ctx context.Context, endpoint string, target interface{}) error {
	fullURL, err := c.buildURL(endpoint)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, http.NoBody)
	if err != nil {
		return err
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, target)
}

func (c *CustomHTTPClient) buildURL(endpoint string) (string, error) {
	base, err := url.Parse(c.BaseURL)
	if err != nil {
		return "", fmt.Errorf("invalid BaseURL: %w", err)
	}
	fullURL := *base
	fullURL.Path = path.Join(fullURL.Path, endpoint)
	return fullURL.String(), nil
}
