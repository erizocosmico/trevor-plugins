package movie

import (
	"github.com/ryanbradynd05/go-tmdb"
	"gopkg.in/mvader/trevor.v1"
	"os"
	"testing"
)

var (
	apiKey = os.Getenv("TMDB_API_KEY")
)

func TestAnalyze(t *testing.T) {
	plugin := &randomMoviePlugin{}

	testCases := []struct {
		text  string
		score float64
		match bool
	}{
		{"", 0, false},
		{"hello world", 0, false},
		{"tell me what movie to watch", 10, true},
		{"tell me what movie should we watch", 10, true},
		{"tell me what to watch", 5, false},
		{"tell me what should we watch", 5, false},
		{"what to watch?", 5, false},
		{"what should we watch?", 5, false},
		{"what movie to watch?", 10, true},
		{"what movie should we watch?", 10, true},
		{"what show should we watch?", 0, false},
		{"tell me a random movie ", 10, true},
		{"random movie", 10, true},
		{"we want to watch something", 5, false},
		{"i want to watch something", 5, false},
		{"i want to watch a movie", 10, true},
	}

	for _, tc := range testCases {
		req := trevor.NewRequest(tc.text, nil)
		score, _ := plugin.Analyze(req)
		if score.IsExactMatch() != tc.match {
			t.Errorf("expected exact match to be %d for text '%s'", tc.match, tc.text)
		}

		if score.Score() != tc.score {
			t.Errorf("expected score %f, %f received for text '%s'", tc.score, score.Score(), tc.text)
		}
	}
}

func TestPoke(t *testing.T) {
	p := newRandom()
	p.Poke()

	// Just for code coverage
	p.PokeEvery()
	p.Name()
	p.Precedence()

	if p.lastID == 0 {
		t.Errorf("expected lastID not to be 0")
	}
}

func TestFetchLastID(t *testing.T) {
	ID, err := fetchLastID(newRandom().tmdbAPI)
	if err != nil || ID == 0 {
		t.Errorf("unexpected error: %s", err.Error())
	}
}

func TestFetchMovie(t *testing.T) {
	movie, err := fetchMovie(newRandom().tmdbAPI, 11)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	if movie.Title != "Star Wars: Episode IV - A New Hope" {
		t.Errorf("expected movie to be Star Wars episode IV, instead is %s", movie.Title)
	}
}

func TestGetRandomNumber(t *testing.T) {
	for i := 0; i < 100000; i++ {
		n := getRandomNumber(1500000)
		if n < 11 || n > 1500000 {
			t.Errorf("invalid random number %d", n)
		}
	}
}

func TestFetchRandomMovie(t *testing.T) {
	movie, err := fetchRandomMovie(newRandom().tmdbAPI, 12, 5)
	assertRandomResult(t, movie, err)
}

func TestProcess(t *testing.T) {
	p := newRandom()
	p.lastID = 12
	data, err := p.Process(nil, nil)
	movie := data.(tmdb.Movie)
	assertRandomResult(t, &movie, err)
}

func assertRandomResult(t *testing.T, movie *tmdb.Movie, err error) {
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	if movie.OriginalTitle != "Star Wars" && movie.OriginalTitle != "FindingNemo" {
		t.Errorf("expected movie to be Star Wars or Finding Nemo, instead is %s", movie.OriginalTitle)
	}
}

func newRandom() *randomMoviePlugin {
	return NewRandomMovieRecommender(apiKey, 5).(*randomMoviePlugin)
}
