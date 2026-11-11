package main

import (
	"strings"
	"unicode"
)

func cleanInput(text string) []string {
	var words []string
	var currentWord strings.Builder
	for _, ch := range text {		
		if unicode.IsSpace(ch) {
			if currentWord.Len() > 0 {
				words = append(words, currentWord.String())
				currentWord.Reset()
			}
		} else {
			currentWord.WriteRune(unicode.ToLower(ch))
		}
	}
	if currentWord.Len() > 0 {
		words = append(words, currentWord.String())
	}
	return words
}


