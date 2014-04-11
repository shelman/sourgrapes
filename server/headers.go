package server

import (
	"math/rand"
)

var (
	firstHeaders = []string{
		"So what's tickling your fancy this evening?",
		"Whatcha in the mood to see?",
		"Pick your poison...",
		"Which of these best describes what you want to watch?",
		"Which of these appeals to you?",
	}
	subsequentHeaders = []string{
		"Let's get more specific...",
		"What else might narrow it down?",
		"Any other cravings?",
	}
)

func getHeader(first bool) string {
	toUse := firstHeaders
	if !first {
		toUse = subsequentHeaders
	}
	return toUse[rand.Int()%len(toUse)]
}
