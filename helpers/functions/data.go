package functions

import (
	"sort"
	"strings"
)

func SortSliceCaseInsensitive(input []string) []string {
	sorted := make([]string, len(input))
	copy(sorted, input)

	sort.Slice(sorted, func(i, j int) bool {
		return strings.ToLower(sorted[i]) < strings.ToLower(sorted[j])
	})

	return sorted
}
