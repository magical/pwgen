package main

import (
	"bytes"
	"expvar"
	"flag"
	"fmt"
	"log"
	"net/http"
	"path"
)

const (
	numShortWords = 12
	numLargeWords = 3
	numWords      = numShortWords + numLargeWords
)

//go:generate go run gendict.go -o dict.go

// Number of bits of entropy generated since the start of the program
var entropy expvar.Int

// https://www.eff.org/files/2016/07/18/eff_large_wordlist.txt
// https://www.eff.org/files/2016/09/08/eff_short_wordlist_1.txt
// https://www.eff.org/files/2016/09/08/eff_short_wordlist_2_0.txt

func main() {
	flag.Parse()

	l, err := listen()
	if err != nil {
		log.Fatal(err)
	}

	expvar.Publish("entropy", &entropy)
	http.HandleFunc("/password", handleGen)
	http.HandleFunc("/password/list/", handleList)
	log.Fatal(http.Serve(l, nil))
}

// handleGen generates passwords
func handleGen(w http.ResponseWriter, req *http.Request) {
	seen := make(map[int]bool, numWords)
	words := make([]string, 0, numWords)
	for len(words) < numShortWords {
		n := rand.Intn(len(dictShort))
		for seen[n] {
			n = rand.Intn(len(words))
		}
		seen[n] = true
		words = append(words, dictShort[n])
	}
	for n := range seen {
		delete(seen, n)
	}
	for i := 0; i < numLargeWords; i++ {
		n := rand.Intn(len(dictLarge))
		for seen[n] {
			n = rand.Intn(len(words))
		}
		seen[n] = true
		words = append(words, dictLarge[n])
	}
	var out bytes.Buffer
	err := tmpl.Execute(&out, &tmplContext{Words: words})
	if err != nil {
		log.Print(err)
		http.Error(w, "internal server error", 500)
		return
	}
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	out.WriteTo(w)
}

// handleList prints a word list to stdout
func handleList(w http.ResponseWriter, req *http.Request) {
	name := path.Base(req.URL.Path)
	var dict []string
	switch name {
	case "short":
		dict = dictShort
	case "large":
		dict = dictLarge
	default:
		http.NotFound(w, req)
		return
	}
	for i, word := range dict {
		fmt.Fprintln(w, i, word)
	}
}
