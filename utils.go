package cronroutine

import "slices"

func getFirstElementGreaterThan(slice []int, value int) int {
	if slice == nil {
		panic("slice is nil")
	}
	for _, element := range slice {
		if element > value {
			return element
		}
	}

	return slice[0]
}

func sliceWithStep(minValue int, maxValue int, step int) []int {
	numbers := make([]int, 0, (maxValue-minValue)/step+1)
	for i := minValue; i <= maxValue; i += step {
		numbers = append(numbers, i)
	}

	return numbers
}

func sortUnique(numbers []int) []int {
	unique := make(map[int]struct{})
	for _, number := range numbers {
		unique[number] = struct{}{}
	}

	sorted := make([]int, 0, len(unique))
	for number := range unique {
		sorted = append(sorted, number)
	}
	slices.Sort(sorted)
	return sorted
}
