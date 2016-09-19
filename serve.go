package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"text/scanner"
	"unicode"
)

const numWords = 12

var words []string

// https://www.eff.org/files/2016/07/18/eff_large_wordlist.txt
// https://www.eff.org/files/2016/09/08/eff_short_wordlist_1.txt
// https://www.eff.org/files/2016/09/08/eff_short_wordlist_2_0.txt

func main() {
	words, err := loadwords("eff_large_wordlist.txt")
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/password", func(w http.ResponseWriter, req *http.Request) {
		seen := make(map[int]bool, numWords)
		for i := 0; i < numWords; i++ {
			n := rand.Intn(len(words))
			for seen[n] {
				n = rand.Intn(len(words))
			}
			seen[n] = true
			fmt.Fprintf(w, "%s\n", words[n])
		}
	})
	log.Fatal(http.ListenAndServe(":8081", nil))
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
