package arr

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestClientGetAddsAPIKeyHeaderAndDecodesResponse(t *testing.T) {
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/api/v3/system/status" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if got := r.Header.Get("X-Api-Key"); got != "secret" {
			t.Fatalf("X-Api-Key = %q", got)
		}
		body, _ := json.Marshal(map[string]string{"version": "4.0.0"})
		return response(http.StatusOK, string(body)), nil
	})

	client := NewClient("http://arr.local/", "secret", &http.Client{Transport: transport})

	var out struct {
		Version string `json:"version"`
	}
	if err := client.Get(context.Background(), "/system/status", nil, &out); err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if out.Version != "4.0.0" {
		t.Fatalf("Version = %q", out.Version)
	}
}

func TestClientGetEncodesQuery(t *testing.T) {
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if got := r.URL.Query().Get("term"); got != "The Expanse" {
			t.Fatalf("term = %q", got)
		}
		return response(http.StatusOK, `[]`), nil
	})

	client := NewClient("http://arr.local", "secret", &http.Client{Transport: transport})

	var out []any
	if err := client.Get(context.Background(), "series/lookup", map[string]string{"term": "The Expanse"}, &out); err != nil {
		t.Fatalf("Get() error = %v", err)
	}
}

func TestClientGetReturnsAPIErrorBody(t *testing.T) {
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return response(http.StatusUnauthorized, "nope\n"), nil
	})

	client := NewClient("http://arr.local", "bad-key", &http.Client{Transport: transport})

	var out map[string]any
	err := client.Get(context.Background(), "health", nil, &out)
	if err == nil {
		t.Fatal("Get() error = nil")
	}
	if got, want := err.Error(), "arr api GET /api/v3/health failed: 401 Unauthorized: nope"; got != want {
		t.Fatalf("error = %q, want %q", got, want)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func response(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Status:     fmt.Sprintf("%d %s", status, http.StatusText(status)),
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
}
