package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func newTestApp(t *testing.T) (*app, http.Handler) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("AppData", t.TempDir())
	a := &app{cfg: config{Langs: []string{"pt-br"}}, client: newClient("", "", ""), token: "secret", host: "127.0.0.1:1"}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/state", a.state)
	mux.HandleFunc("POST /api/config", a.saveConfig)
	mux.HandleFunc("POST /api/watch", a.watch)
	return a, a.guard(mux)
}

func call(h http.Handler, method, path, host, token, body string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, path, strings.NewReader(body))
	r.Host = host
	if token != "" {
		r.Header.Set("X-Token", token)
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w
}

func TestGuard(t *testing.T) {
	_, h := newTestApp(t)
	if w := call(h, "GET", "/api/state", "127.0.0.1:1", "", ""); w.Code != http.StatusForbidden {
		t.Fatalf("no token: %d", w.Code)
	}
	if w := call(h, "GET", "/api/state", "evil.com", "secret", ""); w.Code != http.StatusForbidden {
		t.Fatalf("wrong host: %d", w.Code)
	}
	if w := call(h, "GET", "/api/state", "127.0.0.1:1", "secret", ""); w.Code != http.StatusOK {
		t.Fatalf("valid: %d", w.Code)
	}
}

func TestSaveConfigKeepsPassword(t *testing.T) {
	a, h := newTestApp(t)
	call(h, "POST", "/api/config", a.host, a.token, `{"apiKey":"k","username":"u","password":"p","langs":["en"]}`)
	call(h, "POST", "/api/config", a.host, a.token, `{"apiKey":"k2","username":"u","password":"","langs":["en"]}`)
	if a.cfg.Password != "p" || a.cfg.APIKey != "k2" {
		t.Fatalf("cfg = %+v", a.cfg)
	}
	if w := call(h, "GET", "/api/state", a.host, a.token, ""); strings.Contains(w.Body.String(), `"password":"p"`) {
		t.Fatal("state leaks the password")
	}
	if got := loadConfig(); got.APIKey != "k2" {
		t.Fatalf("saved config = %+v", got)
	}
}

func TestWatchHandler(t *testing.T) {
	a, h := newTestApp(t)
	a.cfg.APIKey = "k"
	if w := call(h, "POST", "/api/watch", a.host, a.token, `{"path":"/does/not/exist"}`); w.Code != http.StatusBadRequest {
		t.Fatalf("bad folder: %d", w.Code)
	}
	dir := t.TempDir()
	call(h, "POST", "/api/watch", a.host, a.token, `{"path":"`+filepath.ToSlash(dir)+`"}`)
	if a.watching == "" || a.stop == nil {
		t.Fatal("not watching")
	}
	call(h, "POST", "/api/watch", a.host, a.token, `{"path":""}`)
	if a.watching != "" || a.stop != nil {
		t.Fatal("still watching")
	}
}

func TestWatchFindsNewStableVideo(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "old.mkv"), []byte("x"), 0o644)
	ctx, cancel := context.WithCancel(context.Background())
	found := make(chan string, 4)
	done := make(chan error)
	go func() { done <- watch(ctx, dir, 20*time.Millisecond, func(v string) { found <- v }) }()

	time.Sleep(30 * time.Millisecond)
	os.WriteFile(filepath.Join(dir, "new.mkv"), []byte("x"), 0o644)
	os.WriteFile(filepath.Join(dir, "has.mkv"), []byte("x"), 0o644)
	os.WriteFile(filepath.Join(dir, "has.srt"), []byte("x"), 0o644)
	select {
	case v := <-found:
		if filepath.Base(v) != "new.mkv" {
			t.Fatalf("found %s", v)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("new video not found")
	}
	time.Sleep(100 * time.Millisecond)
	cancel()
	if err := <-done; err != nil {
		t.Fatal(err)
	}
	if len(found) != 0 {
		t.Fatalf("unexpected extra video %s", <-found)
	}
}
