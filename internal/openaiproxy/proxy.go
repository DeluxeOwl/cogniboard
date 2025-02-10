package openaiproxy

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func NewProxy(openAICompatibleEndpoint string, apiKey string, prefix string) (*httputil.ReverseProxy, error) {
	if openAICompatibleEndpoint == "" {
		return nil, errors.New("empty endpoint")
	}

	if apiKey == "" {
		return nil, errors.New("empty api key")
	}

	url, err := url.Parse(openAICompatibleEndpoint)
	if err != nil {
		return nil, fmt.Errorf("parse url %s: %w", openAICompatibleEndpoint, err)
	}

	proxy := httputil.NewSingleHostReverseProxy(url)
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		// Check if the request has an valid Authorization header.
		// If not or if the provided token is "null", then overwrite it.
		auth := req.Header.Get("Authorization")
		if auth == "" || strings.Contains(strings.ToLower(auth), "null") {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
		}

		// Remove the prefix (e.g "/chat") from the incoming URL path.
		req.URL.Path = strings.TrimPrefix(req.URL.Path, prefix)

		// Ensure that the Host header and URL Host match the target.
		req.Host = url.Host
		req.URL.Host = url.Host
		req.Header.Set("Host", url.Host)
	}

	return proxy, nil
}
