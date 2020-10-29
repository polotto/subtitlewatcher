package messenger

import (
	"encoding/json"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"subtitlewatcher/resources/tmp/locales"
)

func getLocale(defaultLang string, defaultLoc string) (string, string) {
	osHost := runtime.GOOS
	switch osHost {
	case "windows":
		// Exec powershell Get-Culture on Windows.
		cmd := exec.Command("powershell", "Get-Culture | select -exp Name")
		output, err := cmd.Output()
		if err == nil {
			langLocRaw := strings.TrimSpace(string(output))
			langLoc := strings.Split(langLocRaw, "-")
			lang := langLoc[0]
			loc := langLoc[1]
			return lang, loc
		}
	case "darwin":
		cmd := exec.Command("sh", "osascript -e 'user locale of (get system info)'")
		output, err := cmd.Output()
		if err == nil {
			langLocRaw := strings.TrimSpace(string(output))
			langLoc := strings.Split(langLocRaw, "_")
			lang := langLoc[0]
			loc := langLoc[1]
			return lang, loc
		}
	case "linux":
		envLang, ok := os.LookupEnv("LANG")
		if ok {
			langLocRaw := strings.TrimSpace(envLang)
			langLocRaw = strings.Split(envLang, ".")[0]
			langLoc := strings.Split(langLocRaw, "_")
			lang := langLoc[0]
			loc := langLoc[1]
			return lang, loc
		}
	}
	return defaultLang, defaultLoc
}

func ReadMessages() map[string]string {
	defaultLang := "en"
	lang, _ := getLocale(defaultLang, "US")

	var jsonMsg map[string]string
	var jsonBytes []byte
	switch lang {
	case "pb":
		jsonBytes = locales.ResPtJson.StaticContent
	default:
		jsonBytes = locales.ResEnJson.StaticContent
	}

	err := json.Unmarshal(jsonBytes, &jsonMsg)
	if err != nil {
		panic("error to Unmarshal translation json")
	}

	return jsonMsg
}

func Languages() []map[string]string {
	var languages []map[string]string
	err := json.Unmarshal(locales.ResLanguagesJson.StaticContent, &languages)
	if err != nil {
		panic(err)
	}
	return languages
}