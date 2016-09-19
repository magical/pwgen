// +build ignore

package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"text/scanner"
	"unicode"
)

func die(v interface{}) {
	fmt.Fprintln(os.Stderr, v)
	os.Exit(1)
}

func main() {
	output := flag.String("o", "", "output filename")
	flag.Parse()

	dictShort, err := loadwords("eff_short_wordlist_1.txt")
	if err != nil {
		die(err)
	}

	dictLarge, err := loadwords("eff_large_wordlist.txt")
	if err != nil {
		die(err)
	}

	var b bytes.Buffer
	fmt.Fprintln(&b, `// AUTO-GENERATED DO NOT EDIT

package main
`)
	fmt.Fprintf(&b, "var dictShort = []string{\n")
	for _, s := range dictShort {
		fmt.Fprintf(&b, "%#v,\n", s)
	}
	fmt.Fprintf(&b, "}\n")
	fmt.Fprintf(&b, "var dictLarge = []string{\n")
	for _, s := range dictLarge {
		fmt.Fprintf(&b, "%#v,\n", s)
	}
	fmt.Fprintf(&b, "}\n")

	data, err := format.Source(b.Bytes())
	if err != nil {
		die(err)
	}

	if *output != "" {
		err = ioutil.WriteFile(*output, data, 0644)
	} else {
		_, err = os.Stdout.Write(data)
	}
	if err != nil {
		die(err)
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
