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

	const space = rune(' ')
	const urlSpace = "%20"
	for _, v := range str {
		if v == space {
			sb.WriteString(urlSpace)
			continue
		}
		sb.WriteRune(v)
	}

	return sb.String()
}
