package openaiproxy

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
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

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		// TODO: add some logging
		if errors.Is(err, io.EOF) || strings.Contains(err.Error(), "closed") {
			// Connection closed by client, we can ignore this
			return
		}
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(fmt.Sprintf("Proxy error: %v", err)))
	}

	return proxy, nil
}

type FlushWriter struct {
	w      http.ResponseWriter
	status int
}

func (fw *FlushWriter) Header() http.Header {
	return fw.w.Header()
}

func (fw *FlushWriter) Write(p []byte) (n int, err error) {
	if fw.status == 0 {
		fw.WriteHeader(http.StatusOK)
	}
	n, err = fw.w.Write(p)
	if err != nil {
		return n, err
	}
	if f, ok := fw.w.(http.Flusher); ok {
		f.Flush()
	}
	return n, nil
}

func (fw *FlushWriter) WriteHeader(status int) {
	if fw.status == 0 {
		fw.status = status
		fw.w.WriteHeader(status)
	}
}

func NewEchoHandlerWithSSE(proxy *httputil.ReverseProxy) echo.HandlerFunc {
	return func(c echo.Context) error {
		isStreaming := strings.Contains(c.Request().Header.Get("Accept"), "text/event-stream") ||
			c.QueryParam("stream") == "true"

		if isStreaming {
			c.Response().Header().Set("Content-Type", "text/event-stream")
			c.Response().Header().Set("Cache-Control", "no-cache")
			c.Response().Header().Set("Connection", "keep-alive")
			c.Response().Header().Set("Transfer-Encoding", "chunked")
		}

		fw := &FlushWriter{w: c.Response().Writer}
		proxy.ServeHTTP(fw, c.Request())

		return nil
	}
}
