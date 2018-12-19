package arrays

// Cracking the Coding Interview - 1.4 Palindrome Permutation. Given a string, check
// if it is a permutation of a palindrome.
// Time: O(n)
// Space complexity: O(n)
func permPalin(str string) bool {
	if len(str) == 0 || len(str) == 1 {
		return true
	}

	// Create a map of runes and the number of times they occur.
	count := make(map[rune]int)
	for _, v := range str {
		count[v]++
	}

	const even = 0
	const odd = 1
	switch len(str) % 2 {
	case even:
		for _, v := range count {

			// All values must be even to be a permutation of a palindrome.
			if v%2 != even {
				return false
			}
		}
	case odd:
		var countOdd int
		for _, v := range count {

			// All values must be even except one can be odd.
			if v%2 != even {
				countOdd++
				if countOdd > 1 {
					return false
				}
			}
		}
	}
	return true
}
