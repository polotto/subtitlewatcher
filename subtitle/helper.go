package subtitle

import (
	"path/filepath"
	"strings"
)

func GenSubtitleName(inputFile string) string {
	dir, file := filepath.Split(inputFile)
	fileName := strings.Split(file, ".")[0]
	return dir+fileName+".srt"
}

func GenFileName(inputFile string) string {
	_, file := filepath.Split(inputFile)
	fileName := strings.Split(file, ".")[0]
	return fileName
}

func Get(languages []string, inputFile string, errorMsg string) error {
	err := SubtitleDb(languages, inputFile, errorMsg)
	if err == nil {
		return nil
	}
	err = OpenSubtitleDb(languages, inputFile, errorMsg)
	if err == nil {
		return nil
	}
	err = Addic7ed(languages, inputFile, errorMsg)
	if err == nil {
		return nil
	}
	return err
}