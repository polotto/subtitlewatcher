package subtitle

import (
	"fmt"
	"github.com/oz/osdb"
	"log"
)

var osLangs = map[string]string{
	"afr": "afr",
	"alb": "alb",
	"ara": "ara",
	"arm": "arm",
	"bel": "bel",
	"ben": "ben",
	"bos": "bos",
	"bul": "bul",
	"bur": "bur",
	"cat": "cat",
	"chi": "chi",
	"cze": "cze",
	"dan": "dan",
	"dut": "dut",
	"eng": "eng",
	"epo": "epo",
	"est": "est",
	"fin": "fin",
	"fre": "fre",
	"geo": "geo",
	"ger": "ger",
	"glg": "glg",
	"ell": "ell",
	"heb": "heb",
	"hin": "hin",
	"hrv": "hrv",
	"hun": "hun",
	"ice": "ice",
	"ind": "ind",
	"ita": "ita",
	"jpn": "jpn",
	"kaz": "kaz",
	"kor": "kor",
	"lav": "lav",
	"ltz": "ltz",
	"lit": "lit",
	"mac": "mac",
	"mni": "mni",
	"mon": "mon",
	"nor": "nor",
	"per": "per",
	"pol": "pol",
	"por": "por",
	"rus": "rus",
	"scc": "scc",
	"sin": "sin",
	"slo": "slo",
	"slv": "slv",
	"spa": "spa",
	"swa": "swa",
	"swe": "swe",
	"syr": "syr",
	"tam": "tam",
	"tha": "tha",
	"tur": "tur",
	"ukr": "ukr",
	"vie": "vie",
	"rum": "rum",
	"pob": "pob",
	"zht": "zht",
	"zhe": "zhe",
}

func OpenSubtitleDb(languages []string, inputFile string, errorMsg string) error {
	client, err := osdb.NewClient()
	if err != nil {
		return err
	}

	if err = client.LogIn("", "", ""); err != nil {
		return err
	}

	for _, language := range languages {
		subs, err := client.FileSearch(inputFile, []string{osLangs[language]})
		if err != nil {
			return err
		}

		if subs != nil {
			if err := client.DownloadTo(&subs[0], GenSubtitleName(inputFile)); err != nil {
				return err
			}
			log.Println("Original name of subtitle :", subs[0].MovieName)
			break
		}
	}

	return fmt.Errorf(errorMsg)
}
