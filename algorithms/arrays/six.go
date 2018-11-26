package arrays

import (
	"strconv"
	"strings"
)

// Cracking the Coding Interview - 1.6 String Compression. Write a function that does
// basic string compression using the counts of the repeated characters. If the compressed
// string is larger than the original string, then return the original string.
// Time: O(n), where n is the length of the input string.
// Space complexity: O(n)
func compress(str string) string {
	var sb strings.Builder
	runes := []rune(str)
	count := 0
	previousChar := runes[0]

	for _, char := range runes {
		if previousChar == char {
			count++
		} else {
			sb.WriteRune(previousChar)
			sb.WriteString(strconv.Itoa(count))
			previousChar = char
			count = 1
		}
	}

	sb.WriteRune(previousChar)
	sb.WriteString(strconv.Itoa(count))

	compressed := sb.String()
	if len(compressed) < len(str) {
		return compressed
	}

	return compressed
}
