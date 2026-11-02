package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/alecthomas/hasansino/github2karakeep/internal/github"
	"github.com/alecthomas/hasansino/github2karakeep/internal/karakeep"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestExportLimitWarning(t *testing.T) {
	tests := []struct {
		name          string
		repos         []map[string]any
		limit         int
		wantProcessed int
		wantWarning   string
	}{
		{"limited", []map[string]any{repo("one"), repo("two"), repo("three")}, 2, 2, "WARNING: Export limit 2 reached; unprocessed repositories: 1."},
		{"exactly at limit", []map[string]any{repo("one"), repo("two")}, 2, 2, ""},
		{"unlimited", []map[string]any{repo("one"), repo("two"), repo("three")}, 0, 3, ""},
		{"invalid repository before limit", []map[string]any{{"repo": map[string]any{"id": 1}}, repo("one"), repo("two"), repo("three")}, 2, 2, "WARNING: Export limit 2 reached; unprocessed repositories: 1."},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var output bytes.Buffer
			oldOutput, oldTransport := log.Writer(), http.DefaultTransport
			log.SetOutput(&output)
			t.Cleanup(func() { log.SetOutput(oldOutput); http.DefaultTransport = oldTransport })
			var created int
			http.DefaultTransport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
				status, body := http.StatusOK, ""
				switch {
				case r.URL.Host == "api.github.com" && r.URL.Path == "/users/octocat/starred":
					data, err := json.Marshal(tc.repos)
					if err != nil {
						return nil, err
					}
					body = string(data)
				case r.URL.Host == "karakeep.test" && r.Method == http.MethodGet && r.URL.Path == "/api/v1/lists":
					body = `{"lists":[{"id":"list-1","name":"stars"}]}`
				case r.URL.Host == "karakeep.test" && r.Method == http.MethodPost && r.URL.Path == "/api/v1/bookmarks":
					created++
					status, body = http.StatusCreated, fmt.Sprintf(`{"id":"bookmark-%d"}`, created)
				case r.URL.Host == "karakeep.test" && r.Method == http.MethodPut && strings.HasPrefix(r.URL.Path, "/api/v1/lists/list-1/bookmarks/"):
					status = http.StatusNoContent
				default:
					return nil, fmt.Errorf("unexpected request: %s %s", r.Method, r.URL)
				}
				return &http.Response{StatusCode: status, Status: fmt.Sprintf("%d %s", status, http.StatusText(status)), Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body)), Request: r}, nil
			})
			err := run(context.Background(), tc.limit, github.New(time.Second, "test-token", 100), "octocat", karakeep.New(time.Second, "https://karakeep.test", "test-token", ""), "stars", false)
			if err != nil {
				t.Fatal(err)
			}
			if created != tc.wantProcessed {
				t.Errorf("created %d bookmarks, want %d", created, tc.wantProcessed)
			}
			if tc.wantWarning == "" && strings.Contains(output.String(), "Export limit") {
				t.Errorf("unexpected limit warning: %s", output.String())
			}
			if tc.wantWarning != "" && (!strings.Contains(output.String(), tc.wantWarning) || !strings.Contains(output.String(), "Each run starts from the beginning. Set EXPORT_LIMIT=0")) {
				t.Errorf("missing actionable limit warning: %s", output.String())
			}
		})
	}
}

func repo(name string) map[string]any {
	return map[string]any{"repo": map[string]any{"id": 1, "full_name": "owner/" + name, "html_url": "https://github.com/owner/" + name}}
}
