package arrays

import (
	"strings"
)

// Cracking the Coding Interview - 1.3 URLify
// Time:
// Space complexity:
func urlify(str string) string {
	return strings.Replace(str, " ", "%20", -1)
}

func urlify2(str string) string {
	var sb strings.Builder

	for _, v := range str {
		if v == rune(' ') {
			sb.WriteString("%20")
		} else {
			// sb.WriteRune(v)
			sb.WriteString(string(v))
		}
	}
	return sb.String()
}
