package ioutil

import (
	"io/ioutil"
	"os"
	"os/user"
)

func UserHome() string {
	usr, err := user.Current()
	if err != nil {
		return "./"
	}
	return usr.HomeDir
}

func ExecutablePath() string {
	ex, err := os.Executable()
	if err != nil {
		return "./"
	}
	return ex
}

func MakeDirectoryIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.Mkdir(path, os.ModeDir|0755)
	}
	return nil
}

// readLines reads a whole file into memory
// and returns a slice of its lines.
func ReadFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return data, err
}

// writeLines writes the lines to the given file.
func WriteFile(path string, data []byte) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	err = ioutil.WriteFile(path, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
