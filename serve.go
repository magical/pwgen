package main

import (
	"bytes"
	"crypto/subtle"
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

func clearInts(s []int) {
	for i := range s {
		s[i] = 0
	}
}

func clearStrings(s []string) {
	for i := range s {
		s[i] = ""
	}
}

// handleGen generates passwords
func handleGen(w http.ResponseWriter, req *http.Request) {
	rand := newRand()
	seen := make([]int, 0, numWords)
	words := make([]string, 0, numWords)
	for len(words) < numShortWords {
		n := rand.Intn(len(dictShort) - len(seen))
		n = indexWithoutReplacement(n, seen)
		seen = insertSorted(seen, n)
		words = append(words, dictShort[n])
	}
	seen = seen[:0]
	for i := 0; i < numLargeWords; i++ {
		n := rand.Intn(len(dictLarge) - len(seen))
		n = indexWithoutReplacement(n, seen)
		seen = insertSorted(seen, n)
		words = append(words, dictLarge[n])
	}
	defer clearStrings(words[:])
	clearInts(seen[:cap(seen)])
	v := rand.Intn(1e6)
	digits := make([]int, 6)
	for i := range digits {
		digits[i] = v % 10
		v /= 10
	}
	defer clearInts(digits[:])
	var out bytes.Buffer
	err := tmpl.Execute(&out, &tmplContext{Words: words, Digits: digits})
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

// returns the minimum and maximum of two values.
// the time taken by this function is independent of the values.
// Its behavior is undefined if a or b are negative or > 2**31 - 1.
func constantTimeMinMax(a, b int) (min, max int) {
	// note: prior to Go 1.26, ConstantTimeLessOrEq only operated in the range 0..2^31-1.
	// starting in Go 1.26, it is defined as boolToInt(a<=b) with compiler magic to make the operation constant time.
	less := subtle.ConstantTimeLessOrEq(a, b)
	min = subtle.ConstantTimeSelect(less, a, b)
	max = a + b - min
	// or:
	// mask := subtle.ConstantTimeSelect(less, 0, a^b)
	// return a^mask, b^mask
	return min, max
}

// precondition: list is sorted
// postcondition: list is sorted
func insertSorted(list []int, x int) []int {
	for i := range list {
		list[i], x = constantTimeMinMax(list[i], x)
	}
	return append(list, x)
}

// precondition: seen is sorted
func indexWithoutReplacement(idx int, seen []int) int {
	for _, gap := range seen {
		idx += subtle.ConstantTimeLessOrEq(gap, idx)
	}
	return idx
}
