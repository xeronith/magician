package main

import (
	_ "embed"
	"encoding/json"
	"maps"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"sync"
)

// openapiJSON is the OpenAPI 3.1 document served at /openapi.json. It is
// embedded so release binaries work from any directory.
//
//go:embed openapi.json
var openapiJSON []byte

// scalarURL pins the Scalar build that the docs page loads from the CDN in
// the browser. The server itself never fetches it.
const scalarURL = "https://cdn.jsdelivr.net/npm/@scalar/api-reference@1.77.0/esm.js"

// docsPage is the HTML page served at /docs.
const docsPage = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>magician API reference</title>
</head>
<body>
  <div id="app"></div>
  <script type="module">
    import { createApiReference } from "` + scalarURL + `";
    createApiReference('#app', { url: '/openapi.json' });
  </script>
</body>
</html>`

// spell is a stored spell.
type spell struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Effect string `json:"effect"`
}

// spellStore keeps spells in memory behind a mutex. They live for as long as
// the process runs and are reset when it exits. IDs start at 1, increase on
// every creation, and are never reused.
type spellStore struct {
	mu     sync.Mutex
	nextID int64
	spells map[int64]spell
}

func newSpellStore() *spellStore {
	return &spellStore{spells: make(map[int64]spell)}
}

// add stores a spell and returns it with its new ID.
func (s *spellStore) add(name, effect string) spell {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nextID++
	added := spell{ID: s.nextID, Name: name, Effect: effect}
	s.spells[added.ID] = added
	return added
}

// get returns the spell with the given ID.
func (s *spellStore) get(id int64) (spell, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	found, ok := s.spells[id]
	return found, ok
}

// delete removes the spell with the given ID and reports whether it existed.
func (s *spellStore) delete(id int64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.spells[id]
	delete(s.spells, id)
	return ok
}

// list returns every spell in ascending ID order.
func (s *spellStore) list() []spell {
	s.mu.Lock()
	defer s.mu.Unlock()
	spells := make([]spell, 0, len(s.spells))
	for _, id := range slices.Sorted(maps.Keys(s.spells)) {
		spells = append(spells, s.spells[id])
	}
	return spells
}

// newSpellHandler builds the HTTP handler for the spell API and its docs.
func newSpellHandler() http.Handler {
	store := newSpellStore()
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Location", "/docs")
		w.WriteHeader(http.StatusFound)
	})

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	mux.HandleFunc("GET /api/spells", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, store.list())
	})

	mux.HandleFunc("POST /api/spells", func(w http.ResponseWriter, r *http.Request) {
		var requested struct {
			Name   string `json:"name"`
			Effect string `json:"effect"`
		}
		if err := json.NewDecoder(r.Body).Decode(&requested); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON")
			return
		}
		name := strings.TrimSpace(requested.Name)
		effect := strings.TrimSpace(requested.Effect)
		switch {
		case name == "":
			writeError(w, http.StatusBadRequest, "name must not be blank")
		case effect == "":
			writeError(w, http.StatusBadRequest, "effect must not be blank")
		default:
			writeJSON(w, http.StatusCreated, store.add(name, effect))
		}
	})

	mux.HandleFunc("GET /api/spells/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := spellID(r.PathValue("id"))
		if err != nil {
			writeError(w, http.StatusNotFound, "spell not found")
			return
		}
		found, ok := store.get(id)
		if !ok {
			writeError(w, http.StatusNotFound, "spell not found")
			return
		}
		writeJSON(w, http.StatusOK, found)
	})

	mux.HandleFunc("DELETE /api/spells/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := spellID(r.PathValue("id"))
		if err != nil {
			writeError(w, http.StatusNotFound, "spell not found")
			return
		}
		if !store.delete(id) {
			writeError(w, http.StatusNotFound, "spell not found")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})

	mux.HandleFunc("GET /openapi.json", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(openapiJSON)
	})

	mux.HandleFunc("GET /docs", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(docsPage))
	})

	// The remaining patterns catch every method the paths above do not
	// support so that all error bodies stay JSON.
	mux.HandleFunc("/api/spells", disallow(http.MethodGet, http.MethodPost))
	mux.HandleFunc("/api/spells/{id}", disallow(http.MethodGet, http.MethodDelete))
	mux.HandleFunc("/health", disallow(http.MethodGet))
	mux.HandleFunc("/openapi.json", disallow(http.MethodGet))
	mux.HandleFunc("/docs", disallow(http.MethodGet))
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		writeError(w, http.StatusNotFound, "not found")
	})

	return mux
}

// spellID parses a path ID; an unparseable ID simply names no spell.
func spellID(raw string) (int64, error) {
	return strconv.ParseInt(raw, 10, 64)
}

// disallow answers with 405 and the methods the path does support.
func disallow(methods ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Allow", strings.Join(methods, ", "))
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// writeJSON writes body as JSON with the given status code. Encode errors are
// discarded: the status and headers are already sent, so nothing meaningful
// remains to report.
func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

// writeError writes a JSON error object with the given status code.
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
