package pulli

import (
	"reflect"
	"testing"
)

func TestFilter_newFilter(t *testing.T) {
	expected := reflect.ValueOf(&filter{}).Type()
	got := reflect.ValueOf(newFilter([]string{}, FilterModeBlackList)).Type()

	if expected != got {
		t.Fatalf("expected newFilte to return %T, but got: %T", expected, got)
	}
}

func TestFilter_newFilter_emptyFilterMode_defaultsToBlacklist(t *testing.T) {
	f := newFilter([]string{}, "")

	if f.filterMode != FilterModeBlackList {
		t.Fatalf("filterMode was expected to be %v, but was: %v", FilterModeBlackList, f.filterMode)
	}
}

//nolint:funlen
func TestFilter_isAllowed(t *testing.T) {
	testDataSet := map[string]struct {
		expectation    string
		inputs         []string
		filters        []string
		filterMode     string
		expectedResult bool
	}{
		"blacklist, no filters, empty input: is allowed": {
			inputs:         []string{""},
			filters:        nil,
			filterMode:     FilterModeBlackList,
			expectedResult: true,
		},
		"blacklist, input, but no filters: is allowed": {
			inputs:         []string{"/projects/test"},
			filters:        nil,
			filterMode:     FilterModeBlackList,
			expectedResult: true,
		},
		"blacklist, input, but non matching filters: is allowed": {
			inputs:         []string{"/projects/test"},
			filters:        []string{"repositories", "testrepo"},
			filterMode:     FilterModeBlackList,
			expectedResult: true,
		},
		"blacklist, filter exactly matches input: not allowed": {
			inputs:         []string{"/projects/test"},
			filters:        []string{"/projects/test", "non-matching-filter"},
			filterMode:     FilterModeBlackList,
			expectedResult: false,
		},
		"blacklist, regex filter matches input: not allowed": {
			inputs:         []string{"/projects/test"},
			filters:        []string{"non-matching-filter", "^.*test$"},
			filterMode:     FilterModeBlackList,
			expectedResult: false,
		},
		"blacklist, simple string as regex filter matches input: not allowed": {
			inputs:         []string{"/projects/test"},
			filters:        []string{"non-matching-filter", "test"},
			filterMode:     FilterModeBlackList,
			expectedResult: false,
		},

		"whitelist, no filters, empty input: is not allowed": {
			inputs:         []string{""},
			filters:        nil,
			filterMode:     FilterModeWhiteList,
			expectedResult: false,
		},
		"whitelist, input, but no filters: is not allowed": {
			inputs:         []string{"/projects/test"},
			filters:        nil,
			filterMode:     FilterModeWhiteList,
			expectedResult: false,
		},
		"whitelist, input, but non matching filters: is not allowed": {
			inputs:         []string{"/projects/test"},
			filters:        []string{"repositories", "testrepo"},
			filterMode:     FilterModeWhiteList,
			expectedResult: false,
		},
		"whitelist, filter exactly matches input: is allowed": {
			inputs:         []string{"/projects/test"},
			filters:        []string{"/projects/test", "non-matching-filter"},
			filterMode:     FilterModeWhiteList,
			expectedResult: true,
		},
		"whitelist, regex filter matches input: is allowed": {
			inputs:         []string{"/projects/test"},
			filters:        []string{"non-matching-filter", "^.*test$"},
			filterMode:     FilterModeWhiteList,
			expectedResult: true,
		},
		"whitelist, simple string as regex filter matches input: is allowed": {
			inputs:         []string{"/projects/test"},
			filters:        []string{"non-matching-filter", "test"},
			filterMode:     FilterModeWhiteList,
			expectedResult: true,
		},
	}
	for testCaseName, testData := range testDataSet {
		t.Run(testCaseName, func(t *testing.T) {
			for inputIndex, input := range testData.inputs {
				f := newFilter(testData.filters, testData.filterMode)
				result := f.isAllowed(input)

				if testData.expectedResult != result {
					t.Fatalf("inputIndex: %v, expected %v, got: %v", inputIndex, testData.expectedResult, result)
				}
			}
		})
	}
}
