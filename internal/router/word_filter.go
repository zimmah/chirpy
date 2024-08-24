package router

import "strings"

func wordFilter(chirp string) string {
	words := strings.Split(chirp, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, isBadWord := badWords[loweredWord]; isBadWord {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}