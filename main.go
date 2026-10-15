// subtitlewatcher downloads subtitles from OpenSubtitles for a video file or
// for every new video that shows up in a watched folder. Without arguments it
// opens a small web interface in the browser.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const usage = `usage:
  subtitlewatcher                      open the graphical interface
  subtitlewatcher [-lang pt-br,en] get <video>...
  subtitlewatcher [-lang pt-br,en] [-every 10s] watch <folder>

  subtitlewatcher version

Settings come from the interface or from these variables:
  OPENSUBTITLES_USERNAME   optional, raises the daily download limit
  OPENSUBTITLES_PASSWORD   optional
  OPENSUBTITLES_API_KEY    only needed when this build has no key built in
`

// Set at build time: -ldflags "-X main.version=1.2.3 -X main.defaultAPIKey=..."
var (
	version       = "0.0.0-dev"
	defaultAPIKey string
)

var videoExts = map[string]bool{
	".avi": true, ".mkv": true, ".mp4": true, ".m4v": true,
	".mov": true, ".mpg": true, ".mpeg": true, ".wmv": true, ".webm": true,
}

type config struct {
	// APIKey overrides the key built into the binary.
	APIKey   string   `json:"apiKey"`
	Username string   `json:"username"`
	Password string   `json:"password"`
	Langs    []string `json:"langs"`
	Folder   string   `json:"folder"`
}

func (c config) key() string {
	if c.APIKey != "" {
		return c.APIKey
	}
	return defaultAPIKey
}

func configPath() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		dir = "."
	}
	return filepath.Join(dir, "subtitlewatcher", "config.json")
}

// loadConfig reads the saved settings; environment variables win over them.
func loadConfig() config {
	cfg := config{Langs: []string{"pt-br", "en"}}
	if data, err := os.ReadFile(configPath()); err == nil {
		json.Unmarshal(data, &cfg)
	}
	for env, field := range map[string]*string{
		"OPENSUBTITLES_API_KEY":  &cfg.APIKey,
		"OPENSUBTITLES_USERNAME": &cfg.Username,
		"OPENSUBTITLES_PASSWORD": &cfg.Password,
	} {
		if v := os.Getenv(env); v != "" {
			*field = v
		}
	}
	return cfg
}

func saveConfig(cfg config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(configPath()), 0o700); err != nil {
		return err
	}
	return os.WriteFile(configPath(), data, 0o600)
}

func main() {
	log.SetFlags(log.Ltime)
	cfg := loadConfig()
	langFlag := flag.String("lang", strings.Join(cfg.Langs, ","), "subtitle languages in order of preference")
	every := flag.Duration("every", 10*time.Second, "how often watch scans the folder")
	flag.Usage = func() { fmt.Fprint(os.Stderr, usage) }
	flag.Parse()

	if flag.Arg(0) == "version" {
		fmt.Println("subtitlewatcher", version)
		return
	}
	if flag.NArg() == 0 {
		if err := gui(cfg); err != nil {
			log.Fatal(err)
		}
		return
	}
	if flag.NArg() < 2 {
		flag.Usage()
		os.Exit(2)
	}
	if cfg.key() == "" {
		log.Fatal("this build has no API key: set OPENSUBTITLES_API_KEY (https://www.opensubtitles.com/consumers)")
	}
	c := newClient(cfg.key(), cfg.Username, cfg.Password)
	langs := strings.Split(*langFlag, ",")

	switch flag.Arg(0) {
	case "get":
		failed := false
		for _, video := range flag.Args()[1:] {
			if !logFetch(c, video, langs) {
				failed = true
			}
		}
		if failed {
			os.Exit(1)
		}
	case "watch":
		dir := flag.Arg(1)
		log.Printf("watching %s", dir)
		err := watch(context.Background(), dir, *every, func(video string) { logFetch(c, video, langs) })
		if err != nil {
			log.Fatal(err)
		}
	default:
		flag.Usage()
		os.Exit(2)
	}
}

func logFetch(c *client, video string, langs []string) bool {
	lang, remaining, err := c.Fetch(video, langs)
	if err != nil {
		log.Printf("%s: %v", video, err)
		return false
	}
	log.Printf("%s: saved %s subtitle (%d downloads left today)", video, lang, remaining)
	return true
}

// watch polls the folder until ctx ends and calls found once for each video
// that appears after it starts, when its size stops changing and it has no
// .srt yet.
func watch(ctx context.Context, dir string, every time.Duration, found func(video string)) error {
	if info, err := os.Stat(dir); err != nil || !info.IsDir() {
		return fmt.Errorf("%s is not a folder", dir)
	}
	seen := scan(dir)
	pending := map[string]int64{}
	tick := time.NewTicker(every)
	defer tick.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-tick.C:
		}
		for path, size := range scan(dir) {
			if _, ok := seen[path]; ok {
				continue
			}
			if prev, ok := pending[path]; !ok || prev != size {
				pending[path] = size // still being copied or downloaded
				continue
			}
			delete(pending, path)
			seen[path] = size
			if _, err := os.Stat(srtPath(path)); err == nil {
				continue
			}
			found(path)
		}
	}
}

// scan returns the size of every visible video file under dir.
func scan(dir string) map[string]int64 {
	found := map[string]int64{}
	filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if strings.HasPrefix(d.Name(), ".") && path != dir {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if d.IsDir() || !videoExts[strings.ToLower(filepath.Ext(path))] {
			return nil
		}
		if info, err := d.Info(); err == nil && info.Size() > 0 {
			found[path] = info.Size()
		}
		return nil
	})
	return found
}
