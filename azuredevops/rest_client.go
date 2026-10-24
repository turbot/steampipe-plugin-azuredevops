package azuredevops

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	restMaxAttempts = 4
	restMaxWait     = 60 * time.Second
)

// restClient issues REST calls directly rather than through the azure-devops-go-api client, which
// discards response headers and so cannot honor Retry-After or X-RateLimit-Remaining.
type restClient struct {
	http          *http.Client
	baseURL       string
	authorization string
	gate          *throttleGate
	now           func() time.Time
	sleep         func(context.Context, time.Duration) error
}

type restStatusError struct {
	URL        string
	StatusCode int
	Message    string
	// Scope names the PAT scope the request needs, for the permission-denied message.
	Scope string
}

// permissionDeniedPrefix starts every error for a rejected or under-scoped PAT, so callers can
// tell a permission gap from a real failure.
const permissionDeniedPrefix = "azuredevops: permission denied"

func (e *restStatusError) Error() string {
	if e.PermissionDenied() {
		msg := fmt.Sprintf("%s: GET %s returned %d", permissionDeniedPrefix, e.URL, e.StatusCode)
		if e.Scope != "" {
			msg += "; the personal access token needs scope " + e.Scope
		}
		if e.Message != "" {
			msg += ": " + e.Message
		}
		return msg
	}
	if e.Message != "" {
		return fmt.Sprintf("azuredevops: GET %s returned %d: %s", e.URL, e.StatusCode, e.Message)
	}
	return fmt.Sprintf("azuredevops: GET %s returned %d", e.URL, e.StatusCode)
}

// PermissionDenied reports a 401 or 403, a 203 sign-in page, or a redirect to sign-in, which is how
// Azure DevOps answers a PAT that lacks the scope or the access.
func (e *restStatusError) PermissionDenied() bool {
	switch {
	case e.StatusCode == http.StatusUnauthorized, e.StatusCode == http.StatusForbidden, e.StatusCode == http.StatusNonAuthoritativeInfo:
		return true
	case e.StatusCode >= 300 && e.StatusCode < 400:
		return true
	}
	return false
}

// throttleGate holds back requests to one organization after Azure DevOps has asked for a delay.
type throttleGate struct {
	mu    sync.Mutex
	until time.Time
}

func (g *throttleGate) wait(ctx context.Context, now func() time.Time, sleep func(context.Context, time.Duration) error) error {
	g.mu.Lock()
	delay := g.until.Sub(now())
	g.mu.Unlock()
	if delay <= 0 {
		return nil
	}
	return sleep(ctx, delay)
}

func (g *throttleGate) holdUntil(t time.Time) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if t.After(g.until) {
		g.until = t
	}
}

var (
	gatesMu sync.Mutex
	gates   = map[string]*throttleGate{}
)

func gateFor(baseURL string) *throttleGate {
	gatesMu.Lock()
	defer gatesMu.Unlock()
	g, ok := gates[baseURL]
	if !ok {
		g = &throttleGate{}
		gates[baseURL] = g
	}
	return g
}

func sleepContext(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}

// noRedirect stops the client following Azure DevOps' redirect to its sign-in page, which would
// otherwise turn an unauthorized request into a 200 HTML page.
func noRedirect(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }

func newRestClient(baseURL, authorization string) *restClient {
	return &restClient{
		http:          &http.Client{Timeout: 2 * time.Minute, CheckRedirect: noRedirect},
		baseURL:       baseURL,
		authorization: authorization,
		gate:          gateFor(baseURL),
		now:           time.Now,
		sleep:         sleepContext,
	}
}

// getJSON GETs urlStr and decodes the body into v. A 429 or 503 is retried after the Retry-After the
// response names; a success that carries Retry-After, or reports no remaining TSTUs, delays the
// next request to the same organization rather than the current one. scope is named in the error
// when the PAT is refused.
func (c *restClient) getJSON(ctx context.Context, urlStr, scope string, v interface{}) error {
	_, err := c.get(ctx, urlStr, scope, v)
	return err
}

func (c *restClient) get(ctx context.Context, urlStr, scope string, v interface{}) (http.Header, error) {
	for attempt := 1; ; attempt++ {
		if err := c.gate.wait(ctx, c.now, c.sleep); err != nil {
			return nil, err
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", c.authorization)
		req.Header.Set("Accept", "application/json")
		// Without this, a rejected PAT is answered with a 203 sign-in page instead of a 401.
		req.Header.Set("X-TFS-FedAuthRedirect", "Suppress")

		resp, err := c.http.Do(req)
		if err != nil {
			return nil, err
		}
		retryAfter, hasRetryAfter := parseRetryAfter(resp.Header.Get("Retry-After"), c.now())

		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusServiceUnavailable {
			_ = resp.Body.Close()
			if attempt >= restMaxAttempts {
				return nil, &restStatusError{URL: urlStr, StatusCode: resp.StatusCode, Message: "still throttled after retrying"}
			}
			if !hasRetryAfter {
				retryAfter = time.Duration(1<<attempt) * time.Second
			}
			c.gate.holdUntil(c.now().Add(min(retryAfter, restMaxWait)))
			continue
		}

		if hasRetryAfter {
			c.gate.holdUntil(c.now().Add(min(retryAfter, restMaxWait)))
		} else if resp.Header.Get("X-RateLimit-Remaining") == "0" {
			if reset, err := strconv.ParseInt(resp.Header.Get("X-RateLimit-Reset"), 10, 64); err == nil {
				c.gate.holdUntil(minTime(time.Unix(reset, 0), c.now().Add(restMaxWait)))
			}
		}

		defer func() { _ = resp.Body.Close() }()
		if resp.StatusCode >= 300 || resp.StatusCode == http.StatusNonAuthoritativeInfo {
			body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
			var apiErr struct {
				Message string `json:"message"`
			}
			_ = json.Unmarshal(body, &apiErr)
			return nil, &restStatusError{URL: urlStr, StatusCode: resp.StatusCode, Message: apiErr.Message, Scope: scope}
		}
		return resp.Header, json.NewDecoder(resp.Body).Decode(v)
	}
}

// continuationHeader carries the token for the next page on Azure DevOps list APIs.
const continuationHeader = "X-Ms-Continuationtoken"

// getList GETs every page of a list endpoint, following the continuation token header, and calls
// emit per item until it returns false. A page may be a bare JSON array or a {count, value} object.
func getList[T any](ctx context.Context, c *restClient, endpoint string, params url.Values, scope string, emit func(T) bool) error {
	seen := map[string]bool{}
	for {
		var raw json.RawMessage
		header, err := c.get(ctx, endpoint+"?"+params.Encode(), scope, &raw)
		if err != nil {
			return err
		}
		items, err := decodeList[T](raw)
		if err != nil {
			return fmt.Errorf("azuredevops: decoding %s: %w", endpoint, err)
		}
		for _, item := range items {
			if !emit(item) {
				return nil
			}
		}
		token := header.Get(continuationHeader)
		if token == "" || seen[token] {
			return nil
		}
		seen[token] = true
		params.Set("continuationToken", token)
	}
}

func decodeList[T any](raw json.RawMessage) ([]T, error) {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null")) {
		return nil, nil
	}
	var items []T
	if trimmed[0] == '[' {
		err := json.Unmarshal(trimmed, &items)
		return items, err
	}
	var page struct {
		Value []T `json:"value"`
	}
	err := json.Unmarshal(trimmed, &page)
	return page.Value, err
}

// serviceURL returns the base URL of an organization-level Azure DevOps service that lives on its
// own host, such as vsaex, auditservice or advsec. Any other base URL (a test server, Azure DevOps
// Server) is returned unchanged.
func serviceURL(baseURL, service string) string {
	base := strings.TrimRight(baseURL, "/")
	u, err := url.Parse(base)
	if err != nil {
		return base
	}
	host := strings.ToLower(u.Host)
	switch {
	case host == "dev.azure.com":
		u.Host = service + ".dev.azure.com"
		return u.String()
	case strings.HasSuffix(host, ".visualstudio.com"):
		org := strings.TrimSuffix(host, ".visualstudio.com")
		return "https://" + service + ".dev.azure.com/" + org
	}
	return base
}

func projectURL(baseURL, project, path string) string {
	return strings.TrimRight(baseURL, "/") + "/" + url.PathEscape(project) + "/_apis/" + path
}

func isRestNotFound(err error) bool {
	var statusErr *restStatusError
	return errors.As(err, &statusErr) && statusErr.StatusCode == http.StatusNotFound
}

// parseRetryAfter reads either delay-seconds or an HTTP date.
func parseRetryAfter(value string, now time.Time) (time.Duration, bool) {
	if value == "" {
		return 0, false
	}
	if secs, err := strconv.ParseFloat(value, 64); err == nil && secs >= 0 {
		return time.Duration(secs * float64(time.Second)), true
	}
	if t, err := http.ParseTime(value); err == nil {
		return max(t.Sub(now), 0), true
	}
	return 0, false
}

func minTime(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}
