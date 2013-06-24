package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

func train(fn string) map[string]int {
	NWORDS := make(map[string]int)
	p := regexp.MustCompile("[a-z]+")
	if w, err := ioutil.ReadFile(fn); err == nil {
		ws := strings.ToLower(string(w))
		for _, c := range p.FindAllString(ws, -1) { NWORDS[c]++ }
	}
	return NWORDS
}

func edits1(word string, c chan string) {
	alphabet := "abcdefghijklmnopqrstuvwxyz"
	type p struct { a, b string }
	var splits []p
	for i := 0; i < len(word)+1; i++ { splits = append(splits, p{word[:i], word[i:]}) }
	for _, s := range splits {
		if len(s.b) > 0 { c <- s.a + s.b[1:] }
		if len(s.b) > 1 { c <- s.a + string(s.b[1]) + string(s.b[0]) + s.b[2:] }
		for _, abc := range alphabet {
			if len(s.b) > 0 { c <- s.a + string(abc) + s.b[1:] }
			c <- s.a + string(abc) + s.b
		}
	}
}

func edits2(word string, c chan string) {
	ch := make(chan string, 1024*1024)
	go func() { edits1(word, ch); ch <- "" }()
	for e := range ch {
		if e == "" { break }
		edits1(e, c)
	}
}

func best(word string, edits func(string, chan string), model map[string]int) string {
	ch := make(chan string, 1024*1024)
	go func() { edits(word,ch); ch <- "" }()
	mf := 0
	c := ""
	for w := range ch {
		if w == "" { break }
		if f, v := model[w]; v && f > mf { mf, c = f, w }
	}
	return c
}

func correct(word string, model map[string]int) string {
	if _, w := model[word]; w { return word }
	if c:= best(word, edits1, model); c != "" { return c }
	if c:= best(word, edits2, model); c != "" { return c }
	return word
}

func main() {
	NWORDS := train("big.txt")
	fmt.Println(correct("helllo", NWORDS))
}
