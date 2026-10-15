package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

const defaultBaseURL = "https://api.opensubtitles.com/api/v1"

var (
	errNotFound = errors.New("subtitle not found")
	errNotFile  = errors.New("not a video file")
	errWrite    = errors.New("cannot save the subtitle")
)

// apiError keeps the HTTP status so callers can tell a wrong password (401)
// from a used up daily quota (406) or rate limiting (429).
type apiError struct {
	Status int
	Msg    string
}

func (e *apiError) Error() string { return e.Msg }

func apiStatus(err error) int {
	var e *apiError
	if errors.As(err, &e) {
		return e.Status
	}
	return 0
}

// account is what the API reports about the logged user.
type account struct {
	Level     string `json:"level"`
	Allowed   int    `json:"allowed_downloads"`
	Remaining int    `json:"remaining_downloads"`
}

// client talks to the OpenSubtitles REST API (https://opensubtitles.stoplight.io).
type client struct {
	mu       sync.Mutex // one request flow at a time, also respects the rate limit
	baseURL  string
	apiKey   string
	username string
	password string
	token    string
	tokenAt  time.Time
	http     *http.Client
}

func newClient(apiKey, username, password string) *client {
	return &client{
		baseURL:  defaultBaseURL,
		apiKey:   apiKey,
		username: username,
		password: password,
		http:     &http.Client{Timeout: 30 * time.Second},
	}
}

// setCredentials swaps the key and login, dropping the current session.
func (c *client) setCredentials(apiKey, username, password string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.apiKey, c.username, c.password, c.token = apiKey, username, password, ""
	c.baseURL = defaultBaseURL
}

type subtitle struct {
	Attributes struct {
		Language      string `json:"language"`
		DownloadCount int    `json:"download_count"`
		MoviehashOK   bool   `json:"moviehash_match"`
		Files         []struct {
			FileID int `json:"file_id"`
		} `json:"files"`
	} `json:"attributes"`
}

// Fetch finds the best subtitle for the video and saves it next to it as .srt.
// langs is ordered by preference. It returns the language found and how many
// downloads are left today.
func (c *client) Fetch(video string, langs []string) (lang string, remaining int, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if info, err := os.Stat(video); err != nil || !info.Mode().IsRegular() {
		return "", 0, fmt.Errorf("%w: %s", errNotFile, video)
	}
	q := url.Values{}
	if hash, err := movieHash(video); err == nil {
		q.Set("moviehash", hash)
	}
	q.Set("query", strings.ToLower(strings.TrimSuffix(filepath.Base(video), filepath.Ext(video))))
	sorted := append([]string(nil), langs...)
	sort.Strings(sorted)
	q.Set("languages", strings.ToLower(strings.Join(sorted, ",")))

	var res struct {
		Data []subtitle `json:"data"`
	}
	if err := c.do("GET", "/subtitles?"+q.Encode(), nil, &res); err != nil {
		return "", 0, err
	}
	fileID, lang := pick(res.Data, langs)
	if fileID == 0 {
		return "", 0, errNotFound
	}

	var dl struct {
		Link      string `json:"link"`
		Remaining int    `json:"remaining"`
	}
	if err := c.login(); err != nil {
		return "", 0, err
	}
	if err := c.do("POST", "/download", map[string]int{"file_id": fileID}, &dl); err != nil {
		return "", 0, err
	}
	resp, err := c.http.Get(dl.Link)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", 0, &apiError{Status: resp.StatusCode, Msg: "download: " + resp.Status}
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, err
	}
	if err := os.WriteFile(srtPath(video), data, 0o644); err != nil {
		return "", 0, fmt.Errorf("%w: %v", errWrite, err)
	}
	return lang, dl.Remaining, nil
}

// Account logs in again and returns the user's level and daily quota.
func (c *client) Account() (account, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	var res struct {
		Data account `json:"data"`
	}
	if c.username == "" {
		return res.Data, &apiError{Status: http.StatusUnauthorized, Msg: "no username"}
	}
	c.token = ""
	if err := c.login(); err != nil {
		return res.Data, err
	}
	err := c.do("GET", "/infos/user", nil, &res)
	return res.Data, err
}

// pick returns the file of the preferred language, favoring moviehash matches
// and then the most downloaded subtitle.
func pick(subs []subtitle, langs []string) (int, string) {
	for _, lang := range langs {
		var best *subtitle
		for i := range subs {
			s := &subs[i]
			if !strings.EqualFold(s.Attributes.Language, lang) || len(s.Attributes.Files) == 0 {
				continue
			}
			if best == nil ||
				(s.Attributes.MoviehashOK && !best.Attributes.MoviehashOK) ||
				(s.Attributes.MoviehashOK == best.Attributes.MoviehashOK && s.Attributes.DownloadCount > best.Attributes.DownloadCount) {
				best = s
			}
		}
		if best != nil {
			return best.Attributes.Files[0].FileID, lang
		}
	}
	return 0, ""
}

// login is optional: without it the API allows 5 downloads per IP per day.
func (c *client) login() error {
	if c.username == "" || (c.token != "" && time.Since(c.tokenAt) < 12*time.Hour) {
		return nil
	}
	var res struct {
		Token   string `json:"token"`
		BaseURL string `json:"base_url"`
	}
	c.token = ""
	if err := c.do("POST", "/login", map[string]string{"username": c.username, "password": c.password}, &res); err != nil {
		if apiStatus(err) == http.StatusUnauthorized {
			c.username = "" // never retry a wrong password, the API limits logins
		}
		return err
	}
	c.token, c.tokenAt = res.Token, time.Now()
	if res.BaseURL != "" && c.baseURL == defaultBaseURL {
		c.baseURL = "https://" + res.BaseURL + "/api/v1"
	}
	return nil
}

func (c *client) do(method, path string, body, out any) error {
	var r io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		r = bytes.NewReader(b)
	}
	req, err := http.NewRequest(method, c.baseURL+path, r)
	if err != nil {
		return err
	}
	req.Header.Set("Api-Key", c.apiKey)
	req.Header.Set("User-Agent", "subtitlewatcher v"+version)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		var e struct {
			Message string `json:"message"`
		}
		msg := fmt.Sprintf("%s %s: %s", method, strings.SplitN(path, "?", 2)[0], resp.Status)
		if json.Unmarshal(data, &e) == nil && e.Message != "" {
			msg += " (" + e.Message + ")"
		}
		return &apiError{Status: resp.StatusCode, Msg: msg}
	}
	return json.Unmarshal(data, out)
}

// movieHash is the OpenSubtitles hash: file size plus the 64-bit little-endian
// sums of the first and last 64 KiB. Like the reference implementation it only
// accepts 128 KiB to 9 GB; other files are searched by name alone.
func movieHash(path string) (string, error) {
	const chunk = 64 * 1024
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return "", err
	}
	if fi.Size() < 2*chunk || fi.Size() > 9_000_000_000 {
		return "", errors.New("file size out of the hash range")
	}
	hash := uint64(fi.Size())
	buf := make([]byte, chunk)
	for _, off := range []int64{0, fi.Size() - chunk} {
		if _, err := f.ReadAt(buf, off); err != nil {
			return "", err
		}
		for i := 0; i < chunk; i += 8 {
			hash += binary.LittleEndian.Uint64(buf[i:])
		}
	}
	return fmt.Sprintf("%016x", hash), nil
}

func srtPath(video string) string {
	return strings.TrimSuffix(video, filepath.Ext(video)) + ".srt"
}
