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

const numWords = 12

var dict []string

// Number of bits of entropy generated since the start of the program
var entropy expvar.Int

// https://www.eff.org/files/2016/07/18/eff_large_wordlist.txt
// https://www.eff.org/files/2016/09/08/eff_short_wordlist_1.txt
// https://www.eff.org/files/2016/09/08/eff_short_wordlist_2_0.txt

func main() {
	var err error
	dict, err = loadwords("eff_short_wordlist_1.txt")
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
	words := make([]string, numWords)
	for i := range words {
		n := rand.Intn(len(dict))
		for seen[n] {
			n = rand.Intn(len(words))
		}
		seen[n] = true
		words[i] = dict[n]
	}
	var out bytes.Buffer
	tmpl.Execute(&out, &tmplContext{Words: words})
	out.WriteTo(w)
}

// handleList prints a word list to stdout
func handleList(w http.ResponseWriter, req *http.Request) {
	name := path.Base(req.URL.Path)
	_ = name
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
