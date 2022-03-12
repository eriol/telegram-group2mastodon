package utils

import "strings"

// Split text in chunks of almost specified size and append the text in footer.
func SplitTextAtChunk(text string, size int, footer string) []string {
	words := strings.SplitN(text, " ", -1)

	size = size - len(footer)
	chunks := []string{}
	var message string
	for i, word := range words {

		if len(message+" "+word) > size {
			chunks = append(chunks, message+footer)
			message = word
			continue
		}
		if i == 0 {
			message += word
		} else {
			message += " " + word
		}
	}
	chunks = append(chunks, message+footer)

	return chunks
}
