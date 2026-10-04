package main

import (
	"learning/ch5-gordle/gordle"
	"os"
)

const maxAttempts = 6

func main() {
	corpus, err := gordle.ReadCorpus("corpus/english.txt")
	if err != nil {
		panic(err)
	}

	g, err := gordle.New(os.Stdin, corpus, maxAttempts)
	if err != nil {
		panic(err)
	}

	g.Play()
}
