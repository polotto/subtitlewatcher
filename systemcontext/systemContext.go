package systemcontext

import (
	"fmt"
	"runtime"
	"subtitlewatcher/ioutil"
)

func Available() bool {
	osHost := runtime.GOOS
	switch osHost {
	case "windows":
		return true
	}

	return false
}

func Configure() error {
	keyName := "subtitlewatcher"
	rightClickText := "open with Subtitle Watcher"

	executablePath := ioutil.ExecutablePath()
	osHost := runtime.GOOS
	switch osHost {
	case "windows":
		return addRegister(keyName, executablePath, rightClickText)

	default:
		return fmt.Errorf("not available for your system")
	}
}
