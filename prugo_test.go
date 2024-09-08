package prugo

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	defaultHeaders := http.Header{
		"Content-Type": {"application/json"},
	}

	tests := []struct {
		name     string
		baseURL  string
		headers  http.Header
		timeout  time.Duration
		expected struct {
			baseURL string
			headers http.Header
			timeout time.Duration
		}
	}{
		{
			name:    "Default values",
			baseURL: "",
			headers: nil,
			timeout: 0,
			expected: struct {
				baseURL string
				headers http.Header
				timeout time.Duration
			}{
				baseURL: "",
				headers: defaultHeaders,
				timeout: time.Second,
			},
		},
		{
			name:    "Custom baseURL",
			baseURL: "https://example.com",
			headers: nil,
			timeout: 0,
			expected: struct {
				baseURL string
				headers http.Header
				timeout time.Duration
			}{
				baseURL: "https://example.com",
				headers: defaultHeaders,
				timeout: time.Second,
			},
		},
		{
			name:    "Custom headers",
			headers: http.Header{"Authorization": {"Bearer token"}},
			timeout: 2 * time.Second,
			expected: struct {
				baseURL string
				headers http.Header
				timeout time.Duration
			}{
				baseURL: "",
				headers: http.Header{"Authorization": {"Bearer token"}},
				timeout: 2 * time.Second,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(tt.baseURL, tt.headers, tt.timeout)
			got := p.(*prugo)

			if got.baseURL != tt.expected.baseURL {
				t.Errorf("New() baseURL = %v, want %v", got.baseURL, tt.expected.baseURL)
			}

			if got.timeout != tt.expected.timeout {
				t.Errorf("New() timeout = %v, want %v", got.timeout, tt.expected.timeout)
			}

			if !headerEqual(got.headers, tt.expected.headers) {
				t.Errorf("New() headers = %v, want %v", got.headers, tt.expected.headers)
			}
		})
	}
}

// headerEqual checks if two http.Header maps are equal.
func headerEqual(a, b http.Header) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if len(v) != len(b[k]) {
			return false
		}
		for i := range v {
			if v[i] != b[k][i] {
				return false
			}
		}
	}
	return true
}

func TestPrugo_Get(t *testing.T) {
	// Set up a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}
		if r.URL.Path != "/test" {
			t.Errorf("expected URL /test, got %s", r.URL.Path)
		}
		for k, v := range r.Header {
			if k == "Content-Type" && v[0] != "application/json" {
				t.Errorf("expected Content-Type application/json, got %s", v[0])
			}
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer mockServer.Close()

	p := New(mockServer.URL, nil, 0)

	// Test with default headers
	res, err := p.Get("/test", nil)
	if err != nil {
		t.Fatalf("Get() returned error: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 OK, got %d", res.StatusCode)
	}

	// Test with custom headers
	customHeaders := http.Header{"X-Custom-Header": {"value"}}
	config := &Config{headers: customHeaders}
	res, err = p.Get("/test", config)
	if err != nil {
		t.Fatalf("Get() returned error: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 OK, got %d", res.StatusCode)
	}
}

func TestPrugo_GetHeaderMerge(t *testing.T) {
	// Set up a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer token" {
			t.Errorf("expected Authorization header Bearer token, got %s", r.Header.Get("Authorization"))
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer mockServer.Close()

	defaultHeaders := http.Header{"Authorization": {"Bearer token"}}
	p := New(mockServer.URL, defaultHeaders, 0)

	customHeaders := http.Header{"X-Custom-Header": {"value"}}
	config := &Config{headers: customHeaders}
	_, err := p.Get("/test", config)
	if err != nil {
		t.Fatalf("Get() returned error: %v", err)
	}
}

func TestPrugo_GetWithNilConfig(t *testing.T) {
	// Set up a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type header application/json, got %s", r.Header.Get("Content-Type"))
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer mockServer.Close()

	p := New(mockServer.URL, nil, 0)
	_, err := p.Get("/test", nil)
	if err != nil {
		t.Fatalf("Get() returned error: %v", err)
	}
}

// Add tests for Post, Put, Patch, and Delete once implemented
