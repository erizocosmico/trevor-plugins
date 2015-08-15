package translator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/mvader/trevor.v1"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"
)

func runServer(port int) {
	server := trevor.NewServer(trevor.Config{
		Plugins: []trevor.Plugin{trevor.Plugin(NewTranslator(os.Getenv("YANDEX_TRANSLATE_API_KEY")))},
		Port:    port,
	})

	server.Run()

	time.Sleep(5 * time.Millisecond)
}

func makeRequest(text string, port int) []byte {
	jsonStr := []byte(fmt.Sprintf("{\"text\": \"%s\"}", text))
	req, err := http.NewRequest("POST", fmt.Sprintf("http://0.0.0.0:%d/process", port), bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	return body
}

func TestNoMatch(t *testing.T) {
	go runServer(8888)

	body := makeRequest("hello world", 8888)
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		panic(err)
	}

	if !response["error"].(bool) {
		t.Errorf("expected error")
	}
}

func TestTranslate(t *testing.T) {
	go runServer(8889)

	body := makeRequest("how do you say bee in spanish?", 8889)
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		panic(err)
	}

	word := response["data"].(map[string]interface{})["text"].([]interface{})[0].(string)
	if word != "abeja" {
		t.Errorf("expected bee to be translated to abeja")
	}
}

func TestGetWordAndLang(t *testing.T) {
	var testCases = []struct {
		input string
		word  string
		lang  string
		ok    bool
	}{
		{"dskfjslkfjsl", "", "", false},
		{"translate bee to english", "bee", "english", true},
		{"trAnslate bee In Turkish", "bee", "turkish", true},
		{"translate my pony 73 to slovenian", "my pony 73", "slovenian", true},
		{"How do you say lion in swedish?", "lion", "swedish", true},
		{"how would you say zebra in spanish", "zebra", "spanish", true},
		{"how can i say dog in japanese", "dog", "japanese", true},
		{"translate my pony 73 to slovenian", "my pony 73", "slovenian", true},
	}

	for _, tc := range testCases {
		word, lang, ok := getWordAndLang(tc.input)
		if ok != tc.ok {
			t.Errorf("expected okayness to be ", tc.ok)
		}

		if word != tc.word {
			t.Errorf("expected word %s to be %s", word, tc.word)
		}

		if lang != tc.lang {
			t.Errorf("expected lang %s to be %s", lang, tc.lang)
		}
	}
}
