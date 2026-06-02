package arr

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func NewClient(baseURL, apiKey string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		apiKey:     apiKey,
		httpClient: httpClient,
	}
}

func (c *Client) Get(ctx context.Context, path string, query map[string]string, out any) error {
	endpoint, err := c.url(path, query)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-Api-Key", c.apiKey)
	req.Header.Set("Accept", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		body, _ := io.ReadAll(io.LimitReader(res.Body, 4096))
		return fmt.Errorf("arr api GET %s failed: %s: %s", req.URL.RequestURI(), res.Status, strings.TrimSpace(string(body)))
	}

	return json.NewDecoder(res.Body).Decode(out)
}

func (c *Client) url(path string, query map[string]string) (string, error) {
	u, err := url.Parse(c.baseURL + "/api/v3/" + strings.TrimLeft(path, "/"))
	if err != nil {
		return "", err
	}

	values := u.Query()
	for key, value := range query {
		values.Set(key, value)
	}
	u.RawQuery = values.Encode()
	return u.String(), nil
}
