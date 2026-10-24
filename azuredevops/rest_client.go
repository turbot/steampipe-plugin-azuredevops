package azuredevops

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
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
}

func (e *restStatusError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("azuredevops: GET %s returned %d: %s", e.URL, e.StatusCode, e.Message)
	}
	return fmt.Sprintf("azuredevops: GET %s returned %d", e.URL, e.StatusCode)
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

func newRestClient(baseURL, authorization string) *restClient {
	return &restClient{
		http:          &http.Client{Timeout: 2 * time.Minute},
		baseURL:       baseURL,
		authorization: authorization,
		gate:          gateFor(baseURL),
		now:           time.Now,
		sleep:         sleepContext,
	}
}

// getJSON GETs urlStr and decodes the body into v. A 429 or 503 is retried after the Retry-After the
// response names; a success that carries Retry-After, or reports no remaining TSTUs, delays the
// next request to the same organization rather than the current one.
func (c *restClient) getJSON(ctx context.Context, urlStr string, v interface{}) error {
	for attempt := 1; ; attempt++ {
		if err := c.gate.wait(ctx, c.now, c.sleep); err != nil {
			return err
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", c.authorization)
		req.Header.Set("Accept", "application/json")
		// Without this, a rejected PAT is answered with a 203 sign-in page instead of a 401.
		req.Header.Set("X-TFS-FedAuthRedirect", "Suppress")

		resp, err := c.http.Do(req)
		if err != nil {
			return err
		}
		retryAfter, hasRetryAfter := parseRetryAfter(resp.Header.Get("Retry-After"), c.now())

		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusServiceUnavailable {
			_ = resp.Body.Close()
			if attempt >= restMaxAttempts {
				return &restStatusError{URL: urlStr, StatusCode: resp.StatusCode, Message: "still throttled after retrying"}
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
		if resp.StatusCode >= 400 || resp.StatusCode == http.StatusNonAuthoritativeInfo {
			body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
			var apiErr struct {
				Message string `json:"message"`
			}
			_ = json.Unmarshal(body, &apiErr)
			return &restStatusError{URL: urlStr, StatusCode: resp.StatusCode, Message: apiErr.Message}
		}
		return json.NewDecoder(resp.Body).Decode(v)
	}
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
