package subtitle

import (
	"fmt"
	"github.com/matcornic/addic7ed"
	"log"
)

var addic7edLangs = map[string]string{
	"alb": "Albanian",
	"ara": "Arabic",
	"arm": "Armenian",
	"aze": "Azerbaijani",
	"ben": "Bengali",
	"bos": "Bosnian",
	"bul": "Bulgarian",
	"cat": "Català",
	"chi": "Chinese (Simplified)",
	"hrv": "Croatian",
	"cze": "Czech",
	"dan": "Danish",
	"dut": "Dutch",
	"eng": "English",
	"est": "Estonian",
	"baq": "Euskera",
	"fin": "Finnish",
	"fre": "French",
	"frc": "French (Canadian)",
	"glg": "Galego",
	"ger": "German",
	"ell": "Greek",
	"heb": "Hebrew",
	"hin": "Hindi",
	"hun": "Hungarian",
	"ice": "Icelandic",
	"ind": "Indonesian",
	"ita": "Italian",
	"jpn": "Japanese",
	"tlh": "Klingon",
	"kor": "Korean",
	"lav": "Latvian",
	"lit": "Lithuanian",
	"mac": "Macedonian",
	"mal": "Malay",
	"nor": "Norwegian",
	"per": "Persian",
	"pol": "Polish",
	"por": "Portuguese",
	"pob": "Portuguese (Brazilian)",
	"rum": "Romanian",
	"rus": "Russian",
	"scc": "Serbian (Cyrillic)",
	"sin": "Sinhala",
	"slo": "Slovak",
	"slv": "Slovenian",
	"es":  "Spanish",
	"swe": "Swedish",
	"tam": "Tamil",
	"tha": "Thai",
	"tur": "Turkish",
	"ukr": "Ukrainian",
	"vie": "Vietnamese",
	"zht": "Chinese (Traditional)",
}

func Addic7ed(languages []string, inputFile string, errorMsg string) error {
	c := addic7ed.New()

	for _, language := range languages {
		show, subtitle, err := c.SearchBest(GenFileName(inputFile), addic7edLangs[language])
		if err != nil {
			return err
		}

		if show != "" {
			if err := subtitle.DownloadTo(GenSubtitleName(inputFile)); err != nil {
				return err
			}
			log.Println("Original name of subtitle :", show)
			break
		}
	}

	return fmt.Errorf(errorMsg)
}
