# Translator

This package contains all translation-related plugins. All plugins here use the Yandex Translate API.

## translator [![GoDoc](https://godoc.org/github.com/mvader/trevor-plugins/translator?status.svg)](https://godoc.org/github.com/mvader/trevor-plugins/translator)

Translates a word or phrase to another language. The input language is autodetected but the result language must be specified.
It recognises a variety of inputs to identify if the user requests a translation. The following regular expressions are what the plugin will recognise.
* `/^translate (.+) (to|in) (.+)$/`
* `/^how (do you|would you|can i) say (.+) in ([^\?.]+)\??$/`

Valid result languages are:
* spanish
* russian
* english
* portuguese
* turkish
* swedish
* slovenian
* romanian
* norwegian
* dutch
* lithuanian
* korean
* japanese
* german
* french
* finnish
* danish
* chinese

If the expressions are matched the plugin will return an exact match and a score of `10`.
If the expressions are not matched it will return a score of `0`.

### Usage 

```go
package main

import (
  "gopkg.in/mvader/trevor.v1"
  "github.com/mvader/trevor-plugins/translator"
)

func main() {
  server := trevor.NewServer(trevor.Config{
    Plugins: []trevor.Plugin{translator.NewTranslator("YANDEX TRANSLATE API KEY")},
    Port:    8888,
  })

  server.Run()
}
