package sources

import (
	"testing"
	"time"
)

func TestIsToday(t *testing.T) {
	tests := []struct {
		name     string
		date     time.Time
		expected bool
	}{
		{
			name:     "Today",
			date:     time.Now(),
			expected: true,
		},
		{
			name:     "Not Today",
			date:     time.Now().AddDate(0, 0, -1),
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := IsToday(test.date)
			if actual != test.expected {
				t.Errorf("Expected %v but got %v", test.expected, actual)
			}
		})
	}
}
