package translator

import (
	"encoding/json"
	"errors"
	"gopkg.in/mvader/trevor.v1"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type translatorPlugin struct {
	key string
}

var (
	translateToRegExp = regexp.MustCompile(`^translate (.+) (to|in) (.+)$`)
	howDoYouSayRegExp = regexp.MustCompile(`^how (do you|would you|can i) say (.+) in ([^\?.]+)\??$`)

	// translator can only return exact match or no match at all
	matchScore   = trevor.NewScore(10.0, true)
	noMatchScore = trevor.NewScore(0.0, false)

	langCodes = map[string]string{
		"spanish":    "es",
		"russian":    "ru",
		"english":    "en",
		"portuguese": "pt",
		"turkish":    "tk",
		"swedish":    "sv",
		"slovenian":  "sl",
		"romanian":   "ro",
		"norwegian":  "no",
		"dutch":      "nl",
		"lithuanian": "lt",
		"korean":     "ko",
		"japanese":   "ja",
		"german":     "de",
		"french":     "fr",
		"finnish":    "fi",
		"danish":     "da",
		"chinese":    "zh",
	}
)

// NewTranslator creates a new translator plugin instance with the given Yandes Translate API Key.
func NewTranslator(apiKey string) trevor.Plugin {
	return &translatorPlugin{key: apiKey}
}

func (t *translatorPlugin) Analyze(text string) (trevor.Score, interface{}) {
	word, lang, ok := getWordAndLang(text)
	if _, err := getLangCode(lang); ok && err == nil {
		return matchScore, map[string]string{
			"word": word,
			"lang": lang,
		}
	}

	return noMatchScore, nil
}

func (t *translatorPlugin) Process(text string, metadata interface{}) (interface{}, error) {
	if meta, ok := metadata.(map[string]string); ok {
		lang, word := meta["lang"], meta["word"]
		lang, err := getLangCode(lang)
		if err != nil {
			return nil, err
		}

		response, err := t.request(word, lang)
		if err != nil {
			return nil, err
		}

		var result map[string]interface{}
		err = json.Unmarshal(response, &result)
		if err != nil {
			return nil, err
		}

		return result, nil
	}

	return nil, errors.New("can't process text '" + text + "'")
}

func (t *translatorPlugin) Name() string {
	return "translator"
}

func (t *translatorPlugin) Precedence() int {
	return 1
}

func (t *translatorPlugin) request(word, lang string) ([]byte, error) {
	url := "https://translate.yandex.net/api/v1.5/tr.json/translate?key=" + t.key + "&lang=" + lang + "&text=" + url.QueryEscape(word)

	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, err
}

func getLangCode(lang string) (string, error) {
	code, ok := langCodes[lang]
	if !ok {
		return "", errors.New("unknown language " + lang)
	}

	return code, nil
}

func getWordAndLang(text string) (string, string, bool) {
	text = strings.ToLower(text)

	if translateToRegExp.MatchString(text) {
		matches := translateToRegExp.FindStringSubmatch(text)
		return strings.TrimSpace(matches[1]), strings.TrimSpace(matches[3]), true
	} else if howDoYouSayRegExp.MatchString(text) {
		matches := howDoYouSayRegExp.FindStringSubmatch(text)
		return strings.TrimSpace(matches[2]), strings.TrimSpace(matches[3]), true
	}

	return "", "", false
}
