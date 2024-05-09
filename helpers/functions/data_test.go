package functions

import "testing"

func TestSortSliceCaseInsensitive(t *testing.T) {
	tcs := []struct {
		name            string
		content         []string
		expectedContent []string
	}{
		{
			name:            "TestWithLowerCase",
			content:         []string{"banana", "apple", "cherry"},
			expectedContent: []string{"apple", "banana", "cherry"},
		},
		{
			name:            "TestWithUpperCase",
			content:         []string{"BANANA", "APPLE", "CHERRY"},
			expectedContent: []string{"APPLE", "BANANA", "CHERRY"},
		},
		{
			name:            "TestWithMixedCase",
			content:         []string{"Banana", "apple", "Cherry"},
			expectedContent: []string{"apple", "Banana", "Cherry"},
		},
		{
			name:            "TestWithNumbers",
			content:         []string{"1banana", "2apple", "3cherry"},
			expectedContent: []string{"1banana", "2apple", "3cherry"},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			s := SortSliceCaseInsensitive(tc.content)

			for i, v := range s {
				if v != tc.expectedContent[i] {
					t.Errorf("expected %v, got %v", tc.expectedContent[i], v)
				}
			}
		})
	}
}
