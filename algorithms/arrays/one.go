// Cracking the Coding Interview - Chapter 1 Arrays and Strings

package arrays

// 1.1 Is Unique
// Determine if a string has all unique characters. What if you cannot use
// additional data structures?

// Solution 1 (with additional data structure).
// Time complexity: O(n)
// Space complexity: O(n)
func solutionOneA(str string) bool {
	count := map[rune]int{}

	for _, v := range str {
		if count[v] == 1 {
			return false
		}

		count[v] = 1
	}
	return true
}

// Solution 2 (without additional data structure)
// Time: sort O(n^2)
// Space complexity: O(1)
func solutionOneB(str string) bool {
	for ind1, currChar := range str {
		for ind2, anotherChar := range str {
			if currChar == anotherChar && ind2 > ind1 {
				return false
			}
		}
	}

	return true
}
