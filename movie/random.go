package movie

import (
	"errors"
	"github.com/ryanbradynd05/go-tmdb"
	"gopkg.in/mvader/trevor.v1"
	"math/rand"
	"regexp"
	"strings"
	"sync"
	"time"
)

type randomMoviePlugin struct {
	sync.RWMutex
	tmdbAPI    *tmdb.TMDb
	maxRetries int
	precedence int
	lastID     int
}

func NewRandomMovieRecommender(apiKey string, maxRetries int) trevor.Plugin {
	plugin := &randomMoviePlugin{
		tmdbAPI:    tmdb.Init(apiKey),
		maxRetries: maxRetries,
		precedence: 1,
	}

	return plugin
}

type analyzerExpression struct {
	expr         *regexp.Regexp
	isExactMatch func(*regexp.Regexp, string) bool
}

var (
	expressions = []analyzerExpression{
		analyzerExpression{
			regexp.MustCompile(`^(tell me )?what (movie )?(to|should ([a-z]+)) watch\??$`),
			func(expr *regexp.Regexp, text string) bool {
				matches := expr.FindStringSubmatch(text)
				return strings.TrimSpace(matches[2]) == "movie"
			},
		},
		analyzerExpression{
			regexp.MustCompile(`^(tell ([a-z]+) a )?random movie$`),
			func(expr *regexp.Regexp, text string) bool {
				return true
			},
		},
		analyzerExpression{
			regexp.MustCompile(`^([a-z]+) want to watch (a movie|something)$`),
			func(expr *regexp.Regexp, text string) bool {
				matches := expr.FindStringSubmatch(text)
				return strings.TrimSpace(matches[2]) == "a movie"
			},
		},
	}
)

func (p *randomMoviePlugin) Name() string {
	return "random_movie"
}

func (p *randomMoviePlugin) Precedence() int {
	return p.precedence
}

func (p *randomMoviePlugin) Analyze(req *trevor.Request) (trevor.Score, interface{}) {
	text := strings.ToLower(strings.TrimSpace(req.Text))

	for _, e := range expressions {
		if e.expr.MatchString(text) {
			if e.isExactMatch(e.expr, text) {
				return trevor.NewScore(10, true), nil
			} else {
				return trevor.NewScore(5, false), nil
			}
		}
	}

	return trevor.NewScore(0, false), nil
}

func (p *randomMoviePlugin) Process(req *trevor.Request, metadata interface{}) (interface{}, error) {
	movie, err := fetchRandomMovie(p.tmdbAPI, p.lastID, p.maxRetries)
	if err != nil {
		return nil, err
	}

	return *movie, nil
}

func (p *randomMoviePlugin) PokeEvery() time.Duration {
	return 24 * 7 * time.Hour
}

func (p *randomMoviePlugin) Poke() {
	var id int = 0
	var err error = errors.New("dummy error")
	for i := 0; err != nil && i < p.maxRetries; i++ {
		id, err = fetchLastID(p.tmdbAPI)
	}

	p.Lock()
	p.lastID = id
	p.Unlock()
}

func fetchLastID(api *tmdb.TMDb) (int, error) {
	movie, err := api.GetMovieLatest()
	if err != nil {
		return 0, err
	}

	return movie.ID, nil
}

func fetchMovie(api *tmdb.TMDb, ID int) (*tmdb.Movie, error) {
	var options = map[string]string{"append_to_response": "images,videos"}
	return api.GetMovieInfo(ID, options)
}

func getRandomNumber(lastID int) int {
	var n int
	// There are no valid movies on the TMDb api before ID 11
	// https://www.themoviedb.org/talk/519f83fb760ee3572301499b
	for n < 11 {
		rand.Seed(time.Now().UnixNano())
		n = rand.Intn(lastID)
	}

	return n
}

func fetchRandomMovie(api *tmdb.TMDb, lastID, maxRetries int) (*tmdb.Movie, error) {
	var movie *tmdb.Movie
	var err error

	for i := 0; movie == nil && i < maxRetries; i++ {
		ID := getRandomNumber(lastID)
		movie, err = fetchMovie(api, ID)
	}

	return movie, err
}
