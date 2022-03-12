package utils

import "strings"

// Split text in chunks of almost specified size.
func SplitTextAtChunk(text string, size int) []string {
	words := strings.SplitN(text, " ", -1)

	chunks := []string{}
	var message string
	for _, word := range words {

		if len(message+" "+word) > size {
			chunks = append(chunks, message)
			message = word
			continue
		}
		message += " " + word
	}
	chunks = append(chunks, message)

	return chunks
}
