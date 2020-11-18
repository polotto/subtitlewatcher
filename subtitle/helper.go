package subtitle

import (
	"path/filepath"
	"strings"
)

func sliceFileName(file string) string {
	fileSlices := strings.Split(file, ".")
	return strings.Join(fileSlices[0:len(fileSlices)-1], ".")
}

func GenSubtitleName(inputFile string) string {
	dir, file := filepath.Split(inputFile)
	fileName := sliceFileName(file)
	return dir + fileName + ".srt"
}

func GenFileName(inputFile string) string {
	_, file := filepath.Split(inputFile)
	fileName := sliceFileName(file)
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
