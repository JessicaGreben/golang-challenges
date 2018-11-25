package arrays

import "sort"

// Solution
// Time: nlogn
// Space complexity: O(n)
func isPermutation(str1, str2 string) bool {

	if len(str1) != len(str2) {
		return false
	}

	s1 := []rune(str1)
	sort.Sort(sortRunes(s1))

	s2 := []rune(str2)
	sort.Sort(sortRunes(s2))

	if string(s1) == string(s2) {
		return true
	}

	return false
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
