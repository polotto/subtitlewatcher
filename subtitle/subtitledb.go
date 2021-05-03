package subtitle

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var subdbLangs = map[string]string{
	"dut": "nl",
	"eng": "en",
	"fre": "fr",
	"ita": "it",
	"pol": "pl",
	"spa": "es",
	"swe": "sv",
	"tur": "tr",
	"rum": "ro",
	"pob": "pt",
}

//getHashOfVideo gets the hash used by SubDb to identify a video. Absolutely needed either to download or upload subtitles.
//The hash is composed by taking the first and the last 64kb of the video file, putting all together and generating a md5 of the resulting data (128kb).
func getHashOfVideo(filename string) (string, error) {
	readsize := 64 * 1024 // 64kb

	// Open Video
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("Can't open file %v because of : %v ", filename, err.Error())
	}
	defer file.Close()

	// Get stats of file
	fi, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("can't get stats for file %v because of : %v.", filename, err.Error())
	}

	// Fill a buffer with first bytes of file
	bufB := make([]byte, readsize)
	_, err = file.Read(bufB)
	if err != nil {
		return "", fmt.Errorf("can't read content of file %v because of : %v.", filename, err.Error())
	}

	//Fill a buffer with last bytes of file
	bufE := make([]byte, readsize)
	n, err := file.ReadAt(bufE, fi.Size()-int64(len(bufE)))
	if err != nil {
		return "", fmt.Errorf("file is probably too small, can't read content of file %v because of : %v.", filename, err.Error())
	}
	bufE = bufE[:n]

	// Generates MD5 of both bytes chain
	bufB = append(bufB, bufE...)
	hash := fmt.Sprintf("%x", md5.Sum(bufB))

	return hash, nil
}

// Subtitles get the subtitles from the hash of a video
func subtitles(hash string, language string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://api.thesubdb.com/", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", "SubDB/1.0 (Pyrrot/0.1; http://github.com/jrhames/pyrrot-cli)")

	q := req.URL.Query()
	q.Add("action", "download")
	q.Add("hash", hash)
	q.Add("language", language)
	req.URL.RawQuery = q.Encode()

	fmt.Println(req.URL.String())

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return content, nil
}

// SubtitleDb Get get subtitle in the subtitleDb
func SubtitleDb(languages []string, inputFile string, errorMsg string) error {
	hash, err := getHashOfVideo(inputFile)
	if err != nil {
		return err
	}

	fmt.Println(hash)

	for _, language := range languages {
		subtitle, err := subtitles(hash, subdbLangs[language])

		if len(subtitle) != 0 && err == nil {
			err = ioutil.WriteFile(GenSubtitleName(inputFile), subtitle, 0644)
			if err != nil {
				return err
			}
			return nil
		}
	}

	return fmt.Errorf(errorMsg)
}
