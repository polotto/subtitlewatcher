package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestMovieHash(t *testing.T) {
	dir := t.TempDir()
	data := make([]byte, 200*1024)
	data[0] = 1           // first chunk
	data[len(data)-8] = 2 // last chunk
	data[100*1024] = 0xff // middle, ignored
	path := filepath.Join(dir, "a.mkv")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := movieHash(path)
	if err != nil {
		t.Fatal(err)
	}
	if want := "0000000000032003"; got != want { // 200*1024 + 1 + 2
		t.Fatalf("hash = %s, want %s", got, want)
	}
}

func TestPick(t *testing.T) {
	var subs []subtitle
	add := func(lang string, id, downloads int, hashOK bool) {
		var s subtitle
		s.Attributes.Language = lang
		s.Attributes.DownloadCount = downloads
		s.Attributes.MoviehashOK = hashOK
		s.Attributes.Files = append(s.Attributes.Files, struct {
			FileID int `json:"file_id"`
		}{id})
		subs = append(subs, s)
	}
	add("en", 1, 900, true)
	add("pt-BR", 2, 500, false)
	add("pt-BR", 3, 10, true)
	add("pt-BR", 4, 50, true)

	if id, lang := pick(subs, []string{"pt-br", "en"}); id != 4 || lang != "pt-br" {
		t.Fatalf("pick = %d %s, want 4 pt-br", id, lang)
	}
	if id, _ := pick(subs, []string{"fr", "en"}); id != 1 {
		t.Fatalf("fallback pick = %d, want 1", id)
	}
	if id, _ := pick(subs, []string{"fr"}); id != 0 {
		t.Fatalf("missing language pick = %d, want 0", id)
	}
}

func TestFetch(t *testing.T) {
	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/file.srt" && r.Header.Get("Api-Key") != "key" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		switch r.URL.Path {
		case "/subtitles":
			if q := r.URL.Query(); q.Get("languages") != "en,pt-br" || q.Get("query") != "movie.2020" || q.Get("moviehash") == "" {
				t.Errorf("unexpected query %s", r.URL.RawQuery)
			}
			w.Write([]byte(`{"data":[{"attributes":{"language":"pt-BR","moviehash_match":true,"files":[{"file_id":42}]}}]}`))
		case "/download":
			var body map[string]int
			json.NewDecoder(r.Body).Decode(&body)
			if body["file_id"] != 42 {
				t.Errorf("file_id = %d", body["file_id"])
			}
			w.Write([]byte(`{"link":"` + srv.URL + `/file.srt","remaining":4}`))
		case "/file.srt":
			w.Write([]byte("1\n00:00:01,000 --> 00:00:02,000\nOi\n"))
		}
	}))
	defer srv.Close()

	dir := t.TempDir()
	video := filepath.Join(dir, "Movie.2020.mkv")
	os.WriteFile(video, make([]byte, 128*1024), 0o644)

	c := newClient("key", "", "")
	c.baseURL = srv.URL
	lang, remaining, err := c.Fetch(video, []string{"pt-br", "en"})
	if err != nil || lang != "pt-br" || remaining != 4 {
		t.Fatalf("Fetch = %q, %d, %v", lang, remaining, err)
	}
	got, _ := os.ReadFile(filepath.Join(dir, "Movie.2020.srt"))
	if string(got) != "1\n00:00:01,000 --> 00:00:02,000\nOi\n" {
		t.Fatalf("srt = %q", got)
	}
}

func TestAccount(t *testing.T) {
	logins := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/login":
			logins++
			var body map[string]string
			json.NewDecoder(r.Body).Decode(&body)
			if body["password"] != "right" {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"message":"Invalid username/password"}`))
				return
			}
			w.Write([]byte(`{"token":"tok","base_url":"api.opensubtitles.com"}`))
		case "/infos/user":
			if r.Header.Get("Authorization") != "Bearer tok" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.Write([]byte(`{"data":{"level":"Sub leecher","allowed_downloads":20,"remaining_downloads":18}}`))
		}
	}))
	defer srv.Close()

	c := newClient("key", "me", "wrong")
	c.baseURL = srv.URL
	if _, err := c.Account(); errorCode(err) != "loginFailed" {
		t.Fatalf("wrong password: %v", err)
	}
	c.Account()
	if logins != 1 {
		t.Fatalf("wrong password retried: %d logins", logins)
	}

	c.setCredentials("key", "me", "right")
	c.baseURL = srv.URL
	acc, err := c.Account()
	if err != nil || acc.Remaining != 18 || acc.Allowed != 20 {
		t.Fatalf("Account = %+v, %v", acc, err)
	}
}
