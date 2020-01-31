package pulli

import (
	"testing"
)

func TestHelper_truncateString(t *testing.T) {
	testDataSet := map[string]struct {
		input    string
		limit    int
		expected string
	}{
		"empty input - underflow": {
			input:    "",
			limit:    -1,
			expected: "",
		},
		"empty input - truncation (complete)": {
			input:    "",
			limit:    0,
			expected: "",
		},
		"empty input - overflow": {
			input:    "",
			limit:    1,
			expected: "",
		},

		"non empty input - underflow": {
			input:    "hello world",
			limit:    -1,
			expected: "hello world",
		},

		"non empty input - truncation (complete)": {
			input:    "hello world",
			limit:    0,
			expected: "",
		},
		"non empty input - truncation": {
			input:    "hello world",
			limit:    5,
			expected: "hello",
		},
		"non empty input - truncation (nothing truncated)": {
			input:    "hello world",
			limit:    11,
			expected: "hello world",
		},
		"non empty input - overflow": {
			input:    "hello world",
			limit:    12,
			expected: "hello world",
		},
		"non empty input - wider overflow": {
			input:    "hello world",
			limit:    999,
			expected: "hello world",
		},
	}

	for testCaseName, testData := range testDataSet {
		t.Run(testCaseName, func(t *testing.T) {
			result := truncateString(testData.input, testData.limit)
			if testData.expected != result {
				t.Fatalf("expected: %v, got:%v", testData.expected, result)
			}
		})
	}
}
