// Cracking the Coding Interview - Chapter 1 Arrays and Strings

package main

import (
	"sort"
)

func main() {}

// 1.1 Is Unique
// Determine if a string has all unique characters. What if you cannot use
// additional data structures?

// Solution 1 (with additional data structure).
// Time: O(n)
// Space complexity: O(n)
func solutionOneA(str string) (bool, error) {
	count := map[rune]int{}

	for _, v := range str {
		if count[v] == 1 {
			return false, nil
		}

		count[v] = 1
	}
	return true, nil
}

// Solution 2 (without additional data structure)
// Time: sort O(n^2)
// Space complexity: O(1)
func solutionOneB(str string) (bool, error) {
	for ind1, currChar := range str {
		for ind2, anotherChar := range str {
			if currChar == anotherChar && ind2 > ind1 {
				return false, nil
			}
		}
	}

	return true, nil
}

// Solution 3 (with additional data structure)
// Time: sort O(n log n)
// Space complexity: O(len(str))
func solutionOneC(str string) (bool, error) {

	// sort characters.
	s := []rune(str)
	sort.Sort(sortRunes(s))

	// compare neighboring chars for similarity.
	for i := 0; i < len(s)-1; i++ {
		if s[i] == s[i+1] {
			return false, nil
		}
	}

	return true, nil
}

type sortRunes []rune

func (s sortRunes) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s sortRunes) Len() int {
	return len(s)
}
func (s sortRunes) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
