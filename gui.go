package main

import (
	"context"
	"crypto/rand"
	"embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

//go:embed ui.html
var uiFS embed.FS

var errNoPicker = errors.New("noPicker")

type event struct {
	Time      time.Time `json:"time"`
	File      string    `json:"file"`
	Lang      string    `json:"lang,omitempty"`
	Remaining int       `json:"remaining"`
	Error     string    `json:"error,omitempty"`
}

// app is the state behind the web interface.
type app struct {
	mu       sync.Mutex
	cfg      config
	client   *client
	events   []event
	watching string
	stop     context.CancelFunc
	token    string
	host     string
}

// gui serves the interface on localhost and opens it in the browser.
func gui(cfg config) error {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return err
	}
	b := make([]byte, 16)
	rand.Read(b)
	a := &app{cfg: cfg, client: newClient(cfg.key(), cfg.Username, cfg.Password), token: hex.EncodeToString(b), host: ln.Addr().String()}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		page, _ := uiFS.ReadFile("ui.html")
		w.Write(page)
	})
	mux.HandleFunc("GET /api/state", a.state)
	mux.HandleFunc("POST /api/config", a.saveConfig)
	mux.HandleFunc("POST /api/account", a.account)
	mux.HandleFunc("POST /api/pick", a.pick)
	mux.HandleFunc("POST /api/get", a.get)
	mux.HandleFunc("POST /api/watch", a.watch)
	mux.HandleFunc("POST /api/quit", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
		go func() { time.Sleep(200 * time.Millisecond); os.Exit(0) }()
	})

	url := fmt.Sprintf("http://%s/?t=%s", a.host, a.token)
	fmt.Println("subtitlewatcher is running at", url)
	fmt.Println("keep this window open; close it or press Ctrl+C to quit")
	openBrowser(url)
	return http.Serve(ln, a.guard(mux))
}

// guard blocks other sites: the Host must be ours (DNS rebinding) and API calls
// need the session token in a header, which cross-site requests cannot send.
func (a *app) guard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Host != a.host || (strings.HasPrefix(r.URL.Path, "/api/") && r.Header.Get("X-Token") != a.token) {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (a *app) state(w http.ResponseWriter, r *http.Request) {
	a.mu.Lock()
	defer a.mu.Unlock()
	cfg := a.cfg
	hasPassword := cfg.Password != ""
	cfg.Password = ""
	reply(w, map[string]any{
		"config":      cfg,
		"hasPassword": hasPassword,
		"builtInKey":  defaultAPIKey != "",
		"needKey":     a.cfg.key() == "",
		"version":     version,
		"configPath":  configPath(),
		"watching":    a.watching,
		"events":      a.events,
	})
}

func (a *app) saveConfig(w http.ResponseWriter, r *http.Request) {
	var in config
	if json.NewDecoder(r.Body).Decode(&in) != nil {
		replyErr(w, "badRequest")
		return
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if in.Password == "" && in.Username == a.cfg.Username {
		in.Password = a.cfg.Password // empty field keeps the saved password
	}
	in.Folder = a.cfg.Folder
	if err := saveConfig(in); err != nil {
		replyErr(w, err.Error())
		return
	}
	a.cfg = in
	a.client.setCredentials(in.key(), in.Username, in.Password)
	reply(w, nil)
}

func (a *app) account(w http.ResponseWriter, r *http.Request) {
	acc, err := a.client.Account()
	if err != nil {
		log.Printf("account: %v", err)
		replyErr(w, errorCode(err))
		return
	}
	reply(w, acc)
}

// errorCode turns known failures into codes the interface translates.
func errorCode(err error) string {
	switch {
	case errors.Is(err, errNotFound):
		return "notFound"
	case apiStatus(err) == http.StatusUnauthorized:
		return "loginFailed"
	case apiStatus(err) == http.StatusForbidden:
		return "badKey"
	case apiStatus(err) == http.StatusNotAcceptable:
		return "quota"
	case apiStatus(err) == http.StatusTooManyRequests:
		return "tooMany"
	case apiStatus(err) >= 500:
		return "serverError"
	case errors.Is(err, errNotFile):
		return "notFile"
	case errors.Is(err, errWrite):
		return "writeFailed"
	case errors.As(err, new(*url.Error)):
		return "network"
	}
	return "unexpected"
}

func (a *app) pick(w http.ResponseWriter, r *http.Request) {
	var in struct{ Folder bool }
	json.NewDecoder(r.Body).Decode(&in)
	path, err := pickPath(in.Folder)
	if err != nil {
		replyErr(w, err.Error())
		return
	}
	reply(w, map[string]string{"path": path})
}

func (a *app) get(w http.ResponseWriter, r *http.Request) {
	var in struct{ Path string }
	json.NewDecoder(r.Body).Decode(&in)
	path := strings.TrimSpace(in.Path)
	if a.config().key() == "" {
		replyErr(w, "noKey")
		return
	}
	if info, err := os.Stat(path); err != nil || info.IsDir() {
		replyErr(w, "notFile")
		return
	}
	reply(w, a.fetch(path))
}

func (a *app) watch(w http.ResponseWriter, r *http.Request) {
	var in struct{ Path string }
	json.NewDecoder(r.Body).Decode(&in)
	dir := strings.TrimSpace(in.Path)
	if dir != "" && a.config().key() == "" {
		replyErr(w, "noKey")
		return
	}
	if info, err := os.Stat(dir); dir != "" && (err != nil || !info.IsDir()) {
		replyErr(w, "notFolder")
		return
	}

	a.mu.Lock()
	defer a.mu.Unlock()
	if a.stop != nil {
		a.stop()
		a.stop, a.watching = nil, ""
	}
	if dir == "" {
		reply(w, nil)
		return
	}
	ctx, stop := context.WithCancel(context.Background())
	a.stop, a.watching, a.cfg.Folder = stop, dir, dir
	saveConfig(a.cfg)
	go watch(ctx, dir, 10*time.Second, func(video string) { a.fetch(video) })
	reply(w, nil)
}

func (a *app) config() config {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.cfg
}

func (a *app) fetch(video string) event {
	lang, remaining, err := a.client.Fetch(video, a.config().Langs)
	e := event{Time: time.Now(), File: video, Lang: lang, Remaining: remaining}
	if err != nil {
		e.Error = errorCode(err)
		log.Printf("%s: %v", video, err) // full detail stays in the terminal
	}
	a.mu.Lock()
	a.events = append(a.events, e)
	if len(a.events) > 200 {
		a.events = a.events[1:]
	}
	a.mu.Unlock()
	return e
}

func reply(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func replyErr(w http.ResponseWriter, code string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]string{"error": code})
}

func openBrowser(url string) {
	switch runtime.GOOS {
	case "darwin":
		exec.Command("open", url).Start()
	case "windows":
		exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	default:
		exec.Command("xdg-open", url).Start()
	}
}

// pickPath opens the native file or folder dialog. It returns an empty path
// when the user cancels and errNoPicker when the system has no dialog tool.
func pickPath(folder bool) (string, error) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		what := "file"
		if folder {
			what = "folder"
		}
		cmd = exec.Command("osascript", "-e", "activate", "-e", "POSIX path of (choose "+what+")")
	case "windows":
		dialog, prop := "OpenFileDialog", "FileName"
		if folder {
			dialog, prop = "FolderBrowserDialog", "SelectedPath"
		}
		cmd = exec.Command("powershell", "-NoProfile", "-STA", "-Command",
			"Add-Type -AssemblyName System.Windows.Forms; $f = New-Object System.Windows.Forms.Form -Property @{TopMost=$true}; "+
				"$d = New-Object System.Windows.Forms."+dialog+"; if ($d.ShowDialog($f) -eq 'OK') { $d."+prop+" }")
	default:
		if _, err := exec.LookPath("zenity"); err == nil {
			args := []string{"--file-selection"}
			if folder {
				args = append(args, "--directory")
			}
			cmd = exec.Command("zenity", args...)
		} else if _, err := exec.LookPath("kdialog"); err == nil {
			arg := "--getopenfilename"
			if folder {
				arg = "--getexistingdirectory"
			}
			cmd = exec.Command("kdialog", arg)
		} else {
			return "", errNoPicker
		}
	}
	out, err := cmd.Output()
	var exit *exec.ExitError
	if errors.As(err, &exit) {
		return "", nil // canceled
	}
	if err != nil {
		return "", errNoPicker
	}
	if path := strings.TrimSpace(string(out)); path != "" {
		return filepath.Clean(path), nil
	}
	return "", nil
}
