package prugo

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	baseURL := "http://example.com"
	headers := http.Header{
		"Authorization": {"Bearer token"},
	}
	timeout := 2 * time.Second

	p := New(baseURL, headers, timeout)

	if p == nil {
		t.Fatal("Expected non-nil Prugo instance")
	}

	prugoInstance, ok := p.(*prugo)
	if !ok {
		t.Fatalf("Expected *prugo, got %T", p)
	}

	if prugoInstance.baseURL != baseURL {
		t.Errorf("Expected baseURL %s, got %s", baseURL, prugoInstance.baseURL)
	}

	if prugoInstance.headers.Get("Authorization") != "Bearer token" {
		t.Errorf("Expected Authorization header %s, got %s", "Bearer token", prugoInstance.headers.Get("Authorization"))
	}

	if prugoInstance.timeout != timeout {
		t.Errorf("Expected timeout %v, got %v", timeout, prugoInstance.timeout)
	}
}

func TestGet(t *testing.T) {
	// Create a test server to simulate HTTP responses
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected method GET, got %s", r.Method)
		}

		if r.URL.Path != "/test" {
			t.Errorf("Expected URL path /test, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	// Create a Prugo instance using the test server URL
	headers := http.Header{
		"Authorization": {"Bearer token"},
	}
	p := New(ts.URL, headers, 1*time.Second)

	// Call the Get method
	res, err := p.Get("/test", Config{})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, res.StatusCode)
	}
}

func TestGetWithCustomHeaders(t *testing.T) {
	// Create a test server to simulate HTTP responses
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer token" {
			t.Errorf("Expected Authorization header Bearer token, got %s", r.Header.Get("Authorization"))
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type header application/json, got %s", r.Header.Get("Content-Type"))
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	// Create a Prugo instance using the test server URL
	p := New(ts.URL, nil, 1*time.Second)

	// Call the Get method with custom headers
	customHeaders := http.Header{
		"Authorization": {"Bearer token"},
	}
	res, err := p.Get("/test", Config{headers: customHeaders})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, res.StatusCode)
	}
}
