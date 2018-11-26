package arrays

import (
	"math"
)

// Cracking the Coding Interview - 1.5 One Way. There are 3 ways to edit a string:
// insert, remove, replace. Given 2 strings, write a function to check if they are
// one or zero edits away.
// Time: O(n), where n is the length of the longer string.
// Space complexity: O(n), where n is the length of the new strings we create.
func oneAway(str1, str2 string) bool {

	// If the strings are the same, they are zero edits away.
	if str1 == str2 {
		return true
	}

	lenStr1 := len(str1)
	lenStr2 := len(str2)

	// If the length of the strings differ by more than one character,
	// then they are more than one edit away.
	if math.Abs(float64(lenStr1)-float64(lenStr2)) > 1 {
		return false
	}

	// If the string lengths varies by 0 or 1 then they could be one or zero edits away.
	if math.Abs(float64(lenStr1)-float64(lenStr2)) <= 1 {
		return check(str1, str2)
	}

	return false
}

func check(str1, str2 string) bool {
	r1, r2 := []rune(str1), []rune(str2)
	longest, shortest := r1, r2

	if len(r1) < len(r2) {
		longest, shortest = r2, r1
	}

	countDifferChars := 0

	for i, char := range longest {
		if i == len(shortest) || char != shortest[i] {

			// Check if we can replace a character with one edit.
			// If there is 1 character different then you can replace with one edit.
			if len(str1) == len(str2) {
				countDifferChars++
				if countDifferChars > 1 {
					return false
				}
				continue
			}

			// Check if we can remove/insert a character with one edit.
			remove := string(longest[:i]) + string(longest[i+1:])
			insert := string(shortest[:i]) + string(char) + string(shortest[i:])
			if remove == string(shortest) || insert == string(longest) {
				return true
			}

		}
	}

	if countDifferChars == 1 && len(str1) == len(str2) {
		return true
	}

	return false
}
