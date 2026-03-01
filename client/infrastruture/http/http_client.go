package http

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/url"
)

// RequestOption modifies an HTTP request.
type RequestOption func(*http.Request)

// HttpClient is a struct that implements the HttpRequester interface.
type HttpClient struct {
	Client  *http.Client
	BaseURL string
}

// NewHttpClient creates a new HttpClient with a default HTTP client and an optional base URL.
func NewHttpClient(baseURL string) *HttpClient {
	return &HttpClient{
		Client:  &http.Client{},
		BaseURL: baseURL,
	}
}

// buildURL constructs the full URL by combining the base URL with the given path.
func (h *HttpClient) buildURL(path string) (string, error) {
	base, err := url.Parse(h.BaseURL)
	if err != nil {
		return "", err
	}
	fullPath, err := url.Parse(path)
	if err != nil {
		return "", err
	}
	return base.ResolveReference(fullPath).String(), nil
}

// Post sends a POST request to the specified path with the provided body.
func (h *HttpClient) Post(path string, body io.Reader, authToken string) (io.Reader, error) {
	uri, err := h.buildURL(path)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, uri, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	var resp *http.Response
	if authToken != "" {
		resp, err = h.send(req, WithBearerToken(authToken))
	} else {
		resp, err = h.send(req)
	}

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, errors.New("HTTP POST request failed with status: " + resp.Status)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(responseBody), nil
}

// Get sends a GET request to the specified path.
func (h *HttpClient) Get(path, authToken string) (io.Reader, error) {
	uri, err := h.buildURL(path)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}

	var resp *http.Response
	if authToken != "" {
		resp, err = h.send(req, WithBearerToken(authToken))
	} else {
		resp, err = h.send(req)
	}

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, errors.New("HTTP GET request failed with status: " + resp.Status)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(responseBody), nil
}

func (h *HttpClient) send(req *http.Request, opts ...RequestOption) (*http.Response, error) {
	for _, opt := range opts {
		opt(req)
	}

	return h.Client.Do(req)
}

func WithBearerToken(token string) RequestOption {
	return func(req *http.Request) {
		req.Header.Set("Authorization", "Bearer "+token)
	}
}
