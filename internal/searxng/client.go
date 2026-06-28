package searxng

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Client sends requests to a SearXNG instance.
type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
}

// RawSearchResponse contains a SearXNG response body for non-JSON formats.
type RawSearchResponse struct {
	Body        []byte
	ContentType string
}

// NewClient creates a SearXNG API client.
func NewClient(rawBaseURL string, httpClient *http.Client) (*Client, error) {
	parsed, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse base url: %w", err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return nil, fmt.Errorf("base url must use http or https")
	}
	if parsed.Host == "" {
		return nil, fmt.Errorf("base url must include a host")
	}

	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &Client{baseURL: parsed, httpClient: httpClient}, nil
}

func (c *Client) getSearchEndpoint(query, format string) string {
	endpoint := *c.baseURL
	endpoint.Path = strings.TrimRight(endpoint.Path, "/") + "/search"

	values := endpoint.Query()
	values.Set("q", query)
	values.Set("format", format)
	endpoint.RawQuery = values.Encode()

	return endpoint.String()
}

// SearchRaw sends a search request and returns the response body without decoding it.
func (c *Client) SearchRaw(ctx context.Context, query, format string) (*RawSearchResponse, error) {
	format = strings.TrimSpace(format)
	if format == "" {
		return nil, fmt.Errorf("format must not be empty")
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, c.getSearchEndpoint(query, format), nil)
	if err != nil {
		return nil, fmt.Errorf("create search request: %w", err)
	}
	if strings.EqualFold(format, "json") {
		request.Header.Set("Accept", "application/json")
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("send search request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		body, _ := io.ReadAll(io.LimitReader(response.Body, 1024))
		return nil, fmt.Errorf("search request failed: %s: %s", response.Status, strings.TrimSpace(string(body)))
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("read search response: %w", err)
	}

	return &RawSearchResponse{Body: body, ContentType: response.Header.Get("Content-Type")}, nil
}
