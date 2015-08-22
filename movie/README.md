# Movie

This package contains all movie-related plugins. All plugins here use the [TMDb](http://themoviedb.org) API.

## random [![GoDoc](https://godoc.org/github.com/mvader/trevor-plugins/movie?status.svg)](https://godoc.org/github.com/mvader/trevor-plugins/movie)

Recommends a random movie from the TMDb API.
It recognises a variety of inputs to identify if the user requests a movie recommendation. The following regular expressions are what the plugin will recognise.
* `/^(tell me )?what (movie )?(to|should ([a-z]+)) watch\??$/`
* `/^(tell ([a-z]+) a )?random movie$/`
* `/^([a-z]+) want to watch (a movie|something)$/`

If the expressions are matched and the word `movie` is present in the input the plugin will return an exact match and a score of `10`. If the expressions are matched but the word `movie` is not present it will return a score of `5`.
If the expressions are not matched it will return a score of `0`.

### Usage 

```go
package main

import (
  "gopkg.in/mvader/trevor.v1"
  "github.com/mvader/trevor-plugins/movie"
)

func main() {
  var maxRetries int = 5
  server := trevor.NewServer(trevor.Config{
    Plugins: []trevor.Plugin{movie.NewRandomMovieRecommender("TMDB API KEY", maxRetries)},
    Port:    8888,
  })

  server.Run()
}
