// Cracking the Coding Interview - Chapter 1 Arrays and Strings

package arrays

// 1.1 Is Unique - Determine if a string has all unique code points.
// What if you cannot use additional data structures?

// Solution 1 (with additional data structure).
// Time complexity: O(n)
// Space complexity: O(n)
func solutionOneA(str string) bool {
	points := make(map[rune]struct{})
	for _, r := range str {
		if _, found := points[r]; found {
			return false
		}
		points[r] = struct{}{}
	}
	return true
}

// Solution 2 (without additional data structure)
// Time: sort O(n^2)
// Space complexity: O(1)
func solutionOneB(str string) bool {
	for i, r := range str {
		for _, r2 := range str[i+1:] {
			if r == r2 {
				return false
			}
		}
	}
	return true
}
