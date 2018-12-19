package arrays

// Cracking the Coding Interview - 1.8 Zero. Write an algorithm such that in an MxN
// matrix if an element is 0 then its entire row and column are set to 0.
// Time: O(mn)
func zero(matrix [][]int) {
	var rows []int
	columns := make(map[int]struct{})

	for row, rowValues := range matrix {
		for col, value := range rowValues {
			if value == 0 {
				rows = append(rows, row)
				columns[col] = struct{}{}
				break
			}
		}
	}

	for _, row := range rows {
		zeroRow(matrix, row)
	}

	for column := range columns {
		zeroColumn(matrix, column)
	}
}

func zeroRow(matrix [][]int, row int) {
	for i := range matrix[row] {
		matrix[row][i] = 0
	}
}

func zeroColumn(matrix [][]int, column int) {
	for i := range matrix {
		matrix[i][column] = 0
	}
}
