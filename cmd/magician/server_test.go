package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// request sends a request with the given JSON body to handler.
func request(t *testing.T, handler http.Handler, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	var reader io.Reader
	if body != "" {
		reader = strings.NewReader(body)
	}
	req := httptest.NewRequest(method, path, reader)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec
}

// seed creates spells through the API so tests exercise public behavior.
func seed(t *testing.T, handler http.Handler, spells ...spell) {
	t.Helper()
	for _, sp := range spells {
		body, err := json.Marshal(map[string]string{"name": sp.Name, "effect": sp.Effect})
		if err != nil {
			t.Fatalf("marshal seed spell: %v", err)
		}
		rec := request(t, handler, http.MethodPost, "/api/spells", string(body))
		if rec.Code != http.StatusCreated {
			t.Fatalf("seed %q: got status %d, want %d", sp.Name, rec.Code, http.StatusCreated)
		}
	}
}

func TestSpellHandler(t *testing.T) {
	tests := []struct {
		name       string
		seed       []spell
		method     string
		path       string
		body       string
		wantStatus int
		wantBody   string
		wantAllow  string
	}{
		{
			name:       "health reports ok",
			method:     http.MethodGet,
			path:       "/health",
			wantStatus: http.StatusOK,
			wantBody:   `{"status":"ok"}`,
		},
		{
			name:       "post to health is not allowed",
			method:     http.MethodPost,
			path:       "/health",
			wantStatus: http.StatusMethodNotAllowed,
			wantBody:   `{"error":"method not allowed"}`,
			wantAllow:  "GET",
		},
		{
			name:       "empty list is an empty array",
			method:     http.MethodGet,
			path:       "/api/spells",
			wantStatus: http.StatusOK,
			wantBody:   `[]`,
		},
		{
			name:       "list returns seeded spells in id order",
			seed:       []spell{{Name: "fireball", Effect: "burns"}, {Name: "frost", Effect: "freezes"}},
			method:     http.MethodGet,
			path:       "/api/spells",
			wantStatus: http.StatusOK,
			wantBody: `[{"id":1,"name":"fireball","effect":"burns"},` +
				`{"id":2,"name":"frost","effect":"freezes"}]`,
		},
		{
			name:       "unsupported method on the collection is not allowed",
			method:     http.MethodPut,
			path:       "/api/spells",
			wantStatus: http.StatusMethodNotAllowed,
			wantBody:   `{"error":"method not allowed"}`,
			wantAllow:  "GET, POST",
		},
		{
			name:       "get returns a seeded spell",
			seed:       []spell{{Name: "fireball", Effect: "burns"}},
			method:     http.MethodGet,
			path:       "/api/spells/1",
			wantStatus: http.StatusOK,
			wantBody:   `{"id":1,"name":"fireball","effect":"burns"}`,
		},
		{
			name:       "get an unknown spell is not found",
			method:     http.MethodGet,
			path:       "/api/spells/7",
			wantStatus: http.StatusNotFound,
			wantBody:   `{"error":"spell not found"}`,
		},
		{
			name:       "get a non-numeric id is not found",
			method:     http.MethodGet,
			path:       "/api/spells/abc",
			wantStatus: http.StatusNotFound,
			wantBody:   `{"error":"spell not found"}`,
		},
		{
			name:       "post trims and creates a spell",
			method:     http.MethodPost,
			path:       "/api/spells",
			body:       `{"name":" fireball ","effect":" burns "}`,
			wantStatus: http.StatusCreated,
			wantBody:   `{"id":1,"name":"fireball","effect":"burns"}`,
		},
		{
			name:       "post invalid json is a bad request",
			method:     http.MethodPost,
			path:       "/api/spells",
			body:       `{"name":`,
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"error":"invalid JSON"}`,
		},
		{
			name:       "post an empty body is a bad request",
			method:     http.MethodPost,
			path:       "/api/spells",
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"error":"invalid JSON"}`,
		},
		{
			name:       "post a blank name is a bad request",
			method:     http.MethodPost,
			path:       "/api/spells",
			body:       `{"name":"  ","effect":"burns"}`,
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"error":"name must not be blank"}`,
		},
		{
			name:       "post a blank effect is a bad request",
			method:     http.MethodPost,
			path:       "/api/spells",
			body:       `{"name":"fireball","effect":"  "}`,
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"error":"effect must not be blank"}`,
		},
		{
			name:       "unsupported method on a spell is not allowed",
			method:     http.MethodPatch,
			path:       "/api/spells/1",
			wantStatus: http.StatusMethodNotAllowed,
			wantBody:   `{"error":"method not allowed"}`,
			wantAllow:  "GET, DELETE",
		},
		{
			name:       "delete removes a spell",
			seed:       []spell{{Name: "fireball", Effect: "burns"}},
			method:     http.MethodDelete,
			path:       "/api/spells/1",
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "delete an unknown spell is not found",
			method:     http.MethodDelete,
			path:       "/api/spells/7",
			wantStatus: http.StatusNotFound,
			wantBody:   `{"error":"spell not found"}`,
		},
		{
			name:       "a trailing slash is not found",
			method:     http.MethodGet,
			path:       "/api/spells/",
			wantStatus: http.StatusNotFound,
			wantBody:   `{"error":"not found"}`,
		},
		{
			name:       "unknown path is not found",
			method:     http.MethodGet,
			path:       "/api/potions",
			wantStatus: http.StatusNotFound,
			wantBody:   `{"error":"not found"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := newSpellHandler()
			seed(t, handler, tt.seed...)

			rec := request(t, handler, tt.method, tt.path, tt.body)

			if rec.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d (body %q)", rec.Code, tt.wantStatus, rec.Body.String())
			}
			if allow := rec.Header().Get("Allow"); allow != tt.wantAllow {
				t.Errorf("got Allow %q, want %q", allow, tt.wantAllow)
			}
			if got := strings.TrimSpace(rec.Body.String()); got != tt.wantBody {
				t.Errorf("got body %q, want %q", got, tt.wantBody)
			}
			if tt.wantBody != "" {
				if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
					t.Errorf("got Content-Type %q, want application/json", ct)
				}
			}
		})
	}
}

func TestSpellLifecycle(t *testing.T) {
	handler := newSpellHandler()

	rec := request(t, handler, http.MethodPost, "/api/spells", `{"name":"fireball","effect":"burns"}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("first create: got status %d, want %d", rec.Code, http.StatusCreated)
	}
	if got := strings.TrimSpace(rec.Body.String()); got != `{"id":1,"name":"fireball","effect":"burns"}` {
		t.Errorf("first create: got body %q", got)
	}

	rec = request(t, handler, http.MethodPost, "/api/spells", `{"name":"frost","effect":"freezes"}`)
	if got := strings.TrimSpace(rec.Body.String()); got != `{"id":2,"name":"frost","effect":"freezes"}` {
		t.Errorf("second create: got body %q", got)
	}

	rec = request(t, handler, http.MethodDelete, "/api/spells/1", "")
	if rec.Code != http.StatusNoContent {
		t.Fatalf("delete: got status %d, want %d", rec.Code, http.StatusNoContent)
	}

	rec = request(t, handler, http.MethodGet, "/api/spells/1", "")
	if rec.Code != http.StatusNotFound {
		t.Errorf("get after delete: got status %d, want %d", rec.Code, http.StatusNotFound)
	}

	rec = request(t, handler, http.MethodGet, "/api/spells", "")
	if got := strings.TrimSpace(rec.Body.String()); got != `[{"id":2,"name":"frost","effect":"freezes"}]` {
		t.Errorf("list after delete: got body %q", got)
	}

	rec = request(t, handler, http.MethodPost, "/api/spells", `{"name":"light","effect":"illuminates"}`)
	if got := strings.TrimSpace(rec.Body.String()); got != `{"id":3,"name":"light","effect":"illuminates"}` {
		t.Errorf("create after delete: got body %q, want a fresh id", got)
	}
}

func TestRootRedirectsToDocs(t *testing.T) {
	rec := request(t, newSpellHandler(), http.MethodGet, "/", "")
	if rec.Code != http.StatusFound {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusFound)
	}
	if loc := rec.Header().Get("Location"); loc != "/docs" {
		t.Errorf("got Location %q, want /docs", loc)
	}
	if rec.Body.Len() != 0 {
		t.Errorf("got body %q, want none", rec.Body.String())
	}
}

func TestDocsPage(t *testing.T) {
	rec := request(t, newSpellHandler(), http.MethodGet, "/docs", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusOK)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "text/html; charset=utf-8" {
		t.Errorf("got Content-Type %q, want text/html; charset=utf-8", ct)
	}

	page := rec.Body.String()
	for _, want := range []string{
		`<div id="app"></div>`,
		scalarURL,
		`createApiReference('#app', { url: '/openapi.json' })`,
	} {
		if !strings.Contains(page, want) {
			t.Errorf("docs page does not contain %q", want)
		}
	}
}

func TestOpenAPIServerURLIsRelative(t *testing.T) {
	rec := request(t, newSpellHandler(), http.MethodGet, "/openapi.json", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusOK)
	}

	var doc struct {
		Servers []struct {
			URL string `json:"url"`
		} `json:"servers"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &doc); err != nil {
		t.Fatalf("decode document: %v", err)
	}
	if len(doc.Servers) != 1 {
		t.Fatalf("got %d servers, want 1", len(doc.Servers))
	}
	if got := doc.Servers[0].URL; got != "/" {
		t.Errorf("got server URL %q, want the relative / so Scalar calls the serving origin", got)
	}
}

func TestOpenAPIDocument(t *testing.T) {
	rec := request(t, newSpellHandler(), http.MethodGet, "/openapi.json", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusOK)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("got Content-Type %q, want application/json", ct)
	}

	var doc struct {
		OpenAPI string                     `json:"openapi"`
		Paths   map[string]json.RawMessage `json:"paths"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &doc); err != nil {
		t.Fatalf("decode document: %v", err)
	}
	if doc.OpenAPI != "3.1.0" {
		t.Errorf("got openapi %q, want 3.1.0", doc.OpenAPI)
	}

	for _, path := range []string{"/", "/health", "/api/spells", "/api/spells/{id}", "/openapi.json", "/docs"} {
		if _, ok := doc.Paths[path]; !ok {
			t.Errorf("document does not describe %s", path)
		}
	}

	var components struct {
		Components struct {
			Schemas map[string]json.RawMessage `json:"schemas"`
		} `json:"components"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &components); err != nil {
		t.Fatalf("decode components: %v", err)
	}
	for _, name := range []string{"Spell", "NewSpell", "Error"} {
		if _, ok := components.Components.Schemas[name]; !ok {
			t.Errorf("document does not define the %s schema", name)
		}
	}
}
