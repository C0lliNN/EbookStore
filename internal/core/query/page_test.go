package query

import "testing"

func TestPage_Offset(t *testing.T) {
	tests := []struct {
		name     string
		page     Page
		expected int
	}{
		{
			name:     "when page number is 1 and size is 15, then it should return 0",	
			page:     Page{Number: 1, Size: 15},
			expected: 0,
		},
		{
			name:    "when page number is 2 and size is 15, then it should return 15",
			page:    Page{Number: 2, Size: 15},
			expected: 15,
		},
		{
			name:    "when page number is 3 and size is 20, then it should return 40",
			page:    Page{Number: 3, Size: 20},
			expected: 40,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.page.Offset()
			if actual != tc.expected {
				t.Errorf("expected %d, got %d", tc.expected, actual)
			}
		})
	}
}