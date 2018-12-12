package arrays

// Cracking the Coding Interview - 1.4 Palindrome Permutation. Given a string, check
// if it is a permutation of a palindrome.
// Time: O(n)
// Space complexity: O(n)
func permPalin(str string) bool {
	if len(str) == 0 || len(str) == 1 {
		return true
	}

	// Create a hash map that will store the frequency that each character occurs.
	count := make(map[rune]int)

	// Keep track if all the characters occur an even number of times.
	allCharsEven := false

	// Keep track of how many characters occur an odd number of times.
	totalOdd := 0

	for _, v := range str {

		// If the occurence frequency of the current character is odd,
		//
		if count[v]%2 == 1 {
			allCharsEven = true
			totalOdd--

			// If the occurence frequency of the current character is even,
		} else {
			allCharsEven = false
			totalOdd++
		}

		count[v]++
	}

	// If the length of the string is even, then all characters must occur an
	// even number of times to be a palindrome.
	if len(str)%2 == 0 {
		return allCharsEven
	}

	// If the length of the string is odd, then only one character can occur an odd
	// number of times. All other characters must occur an even number of times.
	if totalOdd > 1 {
		return false
	}

	return true
}
