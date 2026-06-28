package searxng

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSearchRaw(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/search" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if got := r.URL.Query().Get("q"); got != "what is searxng" {
			t.Fatalf("q = %q", got)
		}
		if got := r.URL.Query().Get("format"); got != "json" {
			t.Fatalf("format = %q", got)
		}
		if got := r.Header.Get("Accept"); got != "application/json" {
			t.Fatalf("Accept = %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"query":"what is searxng","results":[]}`))
	}))
	defer server.Close()

	client, err := NewClient(server.URL, server.Client())
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	response, err := client.SearchRaw(context.Background(), "what is searxng", "json")
	if err != nil {
		t.Fatalf("SearchRaw() error = %v", err)
	}

	if string(response.Body) != `{"query":"what is searxng","results":[]}` {
		t.Fatalf("body = %q", response.Body)
	}
}

func TestSearchRawPropagatesFormat(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("format"); got != "csv" {
			t.Fatalf("format = %q", got)
		}
		w.Header().Set("Content-Type", "text/csv")
		_, _ = w.Write([]byte("title,url\nSearXNG,https://docs.searxng.org\n"))
	}))
	defer server.Close()

	client, err := NewClient(server.URL, server.Client())
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	response, err := client.SearchRaw(context.Background(), "what is searxng", "csv")
	if err != nil {
		t.Fatalf("SearchRaw() error = %v", err)
	}
	if string(response.Body) != "title,url\nSearXNG,https://docs.searxng.org\n" {
		t.Fatalf("body = %q", response.Body)
	}
}

func TestSearchReturnsStatusError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "bad request", http.StatusBadRequest)
	}))
	defer server.Close()

	client, err := NewClient(server.URL, server.Client())
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	if _, err := client.SearchRaw(context.Background(), "query", "json"); err == nil {
		t.Fatal("SearchRaw() error = nil, want error")
	}
}
