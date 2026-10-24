package azuredevops

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

// Response shapes below follow the Items API reference responses for api-version 7.1.

const testCommit = "23d0bc5b128a10056dc68afece360d8a0fabb014"

type fakeADO struct {
	*httptest.Server
	mu       sync.Mutex
	requests []*http.Request
}

func newFakeADO(t *testing.T, handler func(w http.ResponseWriter, r *http.Request, n int)) *fakeADO {
	f := &fakeADO{}
	f.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f.mu.Lock()
		f.requests = append(f.requests, r)
		n := len(f.requests)
		f.mu.Unlock()
		handler(w, r, n)
	}))
	t.Cleanup(f.Close)
	return f
}

// testClient records the waits it would have made instead of sleeping.
func (f *fakeADO) client(slept *[]time.Duration) *restClient {
	now := time.Unix(1_800_000_000, 0)
	httpClient := f.Client()
	httpClient.CheckRedirect = noRedirect
	c := &restClient{
		http:          httpClient,
		baseURL:       f.URL + "/myorg",
		authorization: "Basic dGVzdA==",
		gate:          &throttleGate{},
	}
	c.now = func() time.Time { return now }
	c.sleep = func(_ context.Context, d time.Duration) error {
		*slept = append(*slept, d)
		now = now.Add(d)
		return nil
	}
	return c
}

func read(t *testing.T, c *restClient, req gitItemRequest) ([]GitItem, error) {
	t.Helper()
	var rows []GitItem
	err := readGitItems(context.Background(), c, req, func(item GitItem) bool {
		rows = append(rows, item)
		return true
	})
	return rows, err
}

func TestListingReadsTheWholeTreeInOneRequest(t *testing.T) {
	f := newFakeADO(t, func(w http.ResponseWriter, r *http.Request, _ int) {
		_, _ = fmt.Fprintf(w, `{"count":3,"value":[
		  {"objectId":"a1","gitObjectType":"tree","commitId":%[1]q,"path":"/","isFolder":true,"url":"u"},
		  {"objectId":"b2","gitObjectType":"tree","commitId":%[1]q,"path":"/src","isFolder":true,"url":"u"},
		  {"objectId":"c3","gitObjectType":"blob","commitId":%[1]q,"path":"/src/package.json","url":"u"}]}`, testCommit)
	})
	var slept []time.Duration
	rows, err := read(t, f.client(&slept), gitItemRequest{Project: "My Project", RepositoryID: "repo-1", Version: "main"})
	if err != nil {
		t.Fatal(err)
	}
	if len(f.requests) != 1 || len(rows) != 3 {
		t.Fatalf("%d requests, %d rows; want 1 and 3", len(f.requests), len(rows))
	}
	r := f.requests[0]
	if r.URL.Path != "/myorg/My Project/_apis/git/repositories/repo-1/items" {
		t.Errorf("path = %q", r.URL.Path)
	}
	q := r.URL.Query()
	for k, want := range map[string]string{
		"scopePath": "/", "recursionLevel": "full", "api-version": "7.1",
		"versionDescriptor.version": "main", "versionDescriptor.versionType": "branch",
	} {
		if got := q.Get(k); got != want {
			t.Errorf("%s = %q, want %q", k, got, want)
		}
	}
	if q.Get("includeContent") != "" {
		t.Error("a listing asked for content")
	}
	if r.Header.Get("Authorization") != "Basic dGVzdA==" || r.Header.Get("X-TFS-FedAuthRedirect") != "Suppress" {
		t.Errorf("headers = %v", r.Header)
	}
	file := rows[2]
	if file.Path != "/src/package.json" || file.GitObjectType != "blob" || file.IsFolder || file.CommitId != testCommit || file.Size != nil || file.Content != nil {
		t.Errorf("file row = %+v", file)
	}
}

func TestAFullHashIsReadAsACommit(t *testing.T) {
	f := newFakeADO(t, func(w http.ResponseWriter, _ *http.Request, _ int) { _, _ = fmt.Fprint(w, `{"value":[]}`) })
	var slept []time.Duration
	if _, err := read(t, f.client(&slept), gitItemRequest{Project: "p", RepositoryID: "r", Version: testCommit, ScopePath: "/src"}); err != nil {
		t.Fatal(err)
	}
	q := f.requests[0].URL.Query()
	if q.Get("versionDescriptor.versionType") != "commit" || q.Get("scopePath") != "/src" {
		t.Errorf("query = %v", q)
	}
}

func TestFileReadReturnsContentAndEchoesThePath(t *testing.T) {
	const content = "{\n  \"name\": \"widget\"\n}\n"
	f := newFakeADO(t, func(w http.ResponseWriter, _ *http.Request, _ int) {
		_, _ = fmt.Fprintf(w, `{"objectId":"c3","gitObjectType":"blob","commitId":%q,"path":"/package.json",
		  "content":%q,"contentMetadata":{"encoding":65001,"contentType":"application/json","isBinary":false}}`, testCommit, content)
	})
	var slept []time.Duration
	rows, err := read(t, f.client(&slept), gitItemRequest{Project: "p", RepositoryID: "r", Path: "package.json"})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("got %d rows", len(rows))
	}
	row := rows[0]
	if row.Path != "package.json" {
		t.Errorf("path = %q, want the qualifier echoed so Postgres keeps the row", row.Path)
	}
	if row.Content == nil || *row.Content != content || row.Size == nil || *row.Size != int64(len(content)) {
		t.Errorf("row = %+v", row)
	}
	q := f.requests[0].URL.Query()
	if q.Get("path") != "package.json" || q.Get("includeContent") != "true" || q.Get("scopePath") != "" {
		t.Errorf("query = %v", q)
	}
}

func TestBinaryAndFolderReadsHaveNoContent(t *testing.T) {
	for name, body := range map[string]string{
		"binary": `{"gitObjectType":"blob","path":"/logo.png","content":"PNG","contentMetadata":{"isBinary":true}}`,
		"nul":    `{"gitObjectType":"blob","path":"/x.bin","content":"a\u0000b"}`,
		"folder": `{"gitObjectType":"tree","path":"/src","isFolder":true}`,
	} {
		t.Run(name, func(t *testing.T) {
			f := newFakeADO(t, func(w http.ResponseWriter, _ *http.Request, _ int) { _, _ = fmt.Fprint(w, body) })
			var slept []time.Duration
			rows, err := read(t, f.client(&slept), gitItemRequest{Project: "p", RepositoryID: "r", Path: "/x"})
			if err != nil {
				t.Fatal(err)
			}
			if rows[0].Content != nil || rows[0].Size != nil {
				t.Errorf("row = %+v", rows[0])
			}
		})
	}
}

func TestMissingItemIsNoRowButUnauthorizedIsAnError(t *testing.T) {
	status := http.StatusNotFound
	f := newFakeADO(t, func(w http.ResponseWriter, _ *http.Request, _ int) {
		w.WriteHeader(status)
		_, _ = fmt.Fprint(w, `{"message":"TF401174: The item 'x' could not be found."}`)
	})
	var slept []time.Duration
	rows, err := read(t, f.client(&slept), gitItemRequest{Project: "p", RepositoryID: "r", Path: "/x"})
	if err != nil || len(rows) != 0 {
		t.Errorf("404: rows=%v err=%v", rows, err)
	}

	status = http.StatusUnauthorized
	if _, err := read(t, f.client(&slept), gitItemRequest{Project: "p", RepositoryID: "r"}); err == nil {
		t.Error("401 returned no error, so a revoked PAT would read as an empty repository")
	}
	status = http.StatusNonAuthoritativeInfo
	if _, err := read(t, f.client(&slept), gitItemRequest{Project: "p", RepositoryID: "r"}); err == nil {
		t.Error("203 sign-in page returned no error")
	}
}

func TestThrottledRequestIsRetriedAfterRetryAfter(t *testing.T) {
	f := newFakeADO(t, func(w http.ResponseWriter, _ *http.Request, n int) {
		if n == 1 {
			w.Header().Set("Retry-After", "7")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		_, _ = fmt.Fprint(w, `{"value":[{"path":"/"}]}`)
	})
	var slept []time.Duration
	rows, err := read(t, f.client(&slept), gitItemRequest{Project: "p", RepositoryID: "r"})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || len(f.requests) != 2 {
		t.Fatalf("rows=%d requests=%d", len(rows), len(f.requests))
	}
	if len(slept) != 1 || slept[0] != 7*time.Second {
		t.Errorf("slept %v, want [7s]", slept)
	}
}

func TestThrottlingGivesUpAfterMaxAttempts(t *testing.T) {
	f := newFakeADO(t, func(w http.ResponseWriter, _ *http.Request, _ int) {
		w.WriteHeader(http.StatusTooManyRequests)
	})
	var slept []time.Duration
	if _, err := read(t, f.client(&slept), gitItemRequest{Project: "p", RepositoryID: "r"}); err == nil {
		t.Fatal("expected an error")
	}
	if len(f.requests) != restMaxAttempts {
		t.Errorf("%d requests, want %d", len(f.requests), restMaxAttempts)
	}
	if slept[0] != 2*time.Second || slept[1] != 4*time.Second {
		t.Errorf("without Retry-After the waits should back off: %v", slept)
	}
}

// A success that is already near the limit delays the next request, not the one that succeeded.
func TestRateLimitHeadersOnSuccessDelayTheNextRequest(t *testing.T) {
	cases := map[string]func(http.Header, time.Time){
		"retry-after": func(h http.Header, _ time.Time) { h.Set("Retry-After", "3") },
		"remaining-zero": func(h http.Header, now time.Time) {
			h.Set("X-RateLimit-Remaining", "0")
			h.Set("X-RateLimit-Reset", fmt.Sprint(now.Add(3*time.Second).Unix()))
		},
	}
	for name, setHeaders := range cases {
		t.Run(name, func(t *testing.T) {
			var c *restClient
			f := newFakeADO(t, func(w http.ResponseWriter, _ *http.Request, n int) {
				if n == 1 {
					setHeaders(w.Header(), c.now())
				}
				_, _ = fmt.Fprint(w, `{"value":[]}`)
			})
			var slept []time.Duration
			c = f.client(&slept)
			if _, err := read(t, c, gitItemRequest{Project: "p", RepositoryID: "r"}); err != nil {
				t.Fatal(err)
			}
			if len(slept) != 0 {
				t.Fatalf("first request slept %v", slept)
			}
			if _, err := read(t, c, gitItemRequest{Project: "p", RepositoryID: "r"}); err != nil {
				t.Fatal(err)
			}
			if len(slept) != 1 || slept[0] != 3*time.Second {
				t.Errorf("second request slept %v, want [3s]", slept)
			}
		})
	}
}

func TestRetryAfterWaitIsCapped(t *testing.T) {
	f := newFakeADO(t, func(w http.ResponseWriter, _ *http.Request, n int) {
		if n == 1 {
			w.Header().Set("Retry-After", "3600")
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		_, _ = fmt.Fprint(w, `{"value":[]}`)
	})
	var slept []time.Duration
	if _, err := read(t, f.client(&slept), gitItemRequest{Project: "p", RepositoryID: "r"}); err != nil {
		t.Fatal(err)
	}
	if slept[0] != restMaxWait {
		t.Errorf("slept %v, want %v", slept[0], restMaxWait)
	}
}

func TestPathAndScopePathTogetherIsAnError(t *testing.T) {
	var slept []time.Duration
	f := newFakeADO(t, func(http.ResponseWriter, *http.Request, int) {})
	if _, err := read(t, f.client(&slept), gitItemRequest{Project: "p", RepositoryID: "r", Path: "/a", ScopePath: "/b"}); err == nil {
		t.Error("expected an error")
	}
	if len(f.requests) != 0 {
		t.Error("made a request")
	}
}

func TestParseRetryAfterAcceptsSecondsAndDates(t *testing.T) {
	now := time.Date(2026, 10, 13, 0, 0, 0, 0, time.UTC)
	if d, ok := parseRetryAfter("2.5", now); !ok || d != 2500*time.Millisecond {
		t.Errorf("seconds: %v %v", d, ok)
	}
	if d, ok := parseRetryAfter(now.Add(10*time.Second).Format(http.TimeFormat), now); !ok || d != 10*time.Second {
		t.Errorf("date: %v %v", d, ok)
	}
	if _, ok := parseRetryAfter("soon", now); ok {
		t.Error("garbage parsed")
	}
}
