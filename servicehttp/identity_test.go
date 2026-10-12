package servicehttp

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPeerPolicyIsExplicitAndKeepsCaller(t *testing.T) {
	token := strings.Repeat("t", 40)
	for _, policy := range []func(string, string, string) bool{nil, func(c, m, p string) bool { return c == "fixture" && m == "GET" && p == "/allowed" }} {
		h := New(Config{Role: "arbitrary", Peers: map[string]string{"fixture": TokenHash(token)}, Authorize: policy}, nil, func(m *http.ServeMux) {
			m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				if Caller(r.Context()) != "fixture" {
					t.Error("caller lost")
				}
				w.WriteHeader(204)
			})
		})
		for _, path := range []string{"/allowed", "/denied", "/readyz", "/metrics"} {
			r := httptest.NewRequest("GET", path, nil)
			r.Header.Set("Authorization", "Bearer "+token)
			w := httptest.NewRecorder()
			h.ServeHTTP(w, r)
			want := 403
			if policy != nil && path == "/allowed" {
				want = 204
			}
			if w.Code != want {
				t.Fatalf("%s: got %d want %d", path, w.Code, want)
			}
		}
	}
}

func TestPublicAssetsCannotExposeOperationalEndpoints(t *testing.T) {
	h := New(Config{UserFacing: true, PublicPath: func(*http.Request) bool { return true }}, nil, func(m *http.ServeMux) {
		m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(204) })
	})
	for _, path := range []string{"/", "/readyz", "/metrics"} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest("GET", path, nil))
		want := 401
		if path == "/" {
			want = 204
		}
		if w.Code != want {
			t.Fatal(path, w.Code)
		}
	}
}

func TestDuplicateOrMalformedPeerIdentityRejected(t *testing.T) {
	hash := TokenHash(strings.Repeat("t", 40))
	for _, peers := range []map[string]string{{"a": hash, "b": hash}, {"a": "invalid"}, {"": hash}, {"a": strings.ToUpper(hash)}} {
		if ValidatePeers(peers) == nil {
			t.Fatal("accepted ambiguous or invalid identity")
		}
	}
	if err := ValidatePeers(map[string]string{"a": hash}); err != nil {
		t.Fatal(err)
	}
}
