package main

import (
	"bytes"
	"expvar"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"text/scanner"
	"unicode"
)

const (
	numShortWords = 12
	numLargeWords = 3
	numWords      = numShortWords + numLargeWords
)

var (
	dictShort []string
	dictLarge []string
)

// Number of bits of entropy generated since the start of the program
var entropy expvar.Int

// https://www.eff.org/files/2016/07/18/eff_large_wordlist.txt
// https://www.eff.org/files/2016/09/08/eff_short_wordlist_1.txt
// https://www.eff.org/files/2016/09/08/eff_short_wordlist_2_0.txt

func main() {
	var err error
	dictShort, err = loadwords("eff_short_wordlist_1.txt")
	if err != nil {
		log.Fatal(err)
	}
	dictLarge, err = loadwords("eff_large_wordlist.txt")
	if err != nil {
		log.Fatal(err)
	}
	expvar.Publish("entropy", &entropy)
	http.HandleFunc("/password", handleGen)
	http.HandleFunc("/password/list/", handleList)
	log.Fatal(http.ListenAndServe(":8081", nil))
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

func loadwords(filename string) ([]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	var words []string
	var s scanner.Scanner
	s.Init(f)
	s.Whitespace = 1<<'\t' | 1<<'\r' | 1<<' '
	s.Mode = scanner.ScanIdents | scanner.ScanInts | scanner.SkipComments
	s.Position.Filename = filename
	s.IsIdentRune = func(r rune, i int) bool {
		return unicode.IsLetter(r) || r == '-'
	}
	for {
		tok := s.Scan()
		if tok == scanner.EOF {
			break
		}
		if tok == '\n' {
			continue
		}
		if tok != scanner.Int {
			return nil, fmt.Errorf("%s: expected int", s.Pos())
		}
		if tok := s.Scan(); tok != scanner.Ident {
			return nil, fmt.Errorf("%s: expected word, found %v", s.Pos(), scanner.TokenString(tok))
		}
		words = append(words, s.TokenText())
		if tok := s.Scan(); tok != '\n' {
			return nil, fmt.Errorf("%s: expected newline, found %v", s.Pos(), scanner.TokenString(tok))
		}
	}
	return words, nil
}
