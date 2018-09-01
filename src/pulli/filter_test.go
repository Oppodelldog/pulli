package pulli

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFilter_newFilter(t *testing.T) {
	assert.IsType(t, &filter{}, newFilter([]string{}, filterModeBlackList))
}

func TestFilter_newFilter_emptyFilterMode_defaultsToBlacklist(t *testing.T) {
	f := newFilter([]string{}, "")

	assert.Exactly(t, filterModeBlackList, f.filterMode)
}

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
			filterMode:     filterModeBlackList,
			expectedResult: true,
		},
		"blacklist, input, but no filters: is allowed": {
			inputs:         []string{"/projects/test"},
			filters:        nil,
			filterMode:     filterModeBlackList,
			expectedResult: true,
		},
		"blacklist, input, but non matching filters: is allowed": {
			inputs:         []string{"/projects/test"},
			filters:        []string{"repositories", "testrepo"},
			filterMode:     filterModeBlackList,
			expectedResult: true,
		},
		"blacklist, filter exactly matches input: not allowed": {
			inputs:         []string{"/projects/test"},
			filters:        []string{"/projects/test", "non-matching-filter"},
			filterMode:     filterModeBlackList,
			expectedResult: false,
		},
		"blacklist, regex filter matches input: not allowed": {
			inputs:         []string{"/projects/test"},
			filters:        []string{"non-matching-filter", "^.*test$"},
			filterMode:     filterModeBlackList,
			expectedResult: false,
		},
		"blacklist, simple string as regex filter matches input: not allowed": {
			inputs:         []string{"/projects/test"},
			filters:        []string{"non-matching-filter", "test"},
			filterMode:     filterModeBlackList,
			expectedResult: false,
		},

		"whitelist, no filters, empty input: is not allowed": {
			inputs:         []string{""},
			filters:        nil,
			filterMode:     filterModeWhiteList,
			expectedResult: false,
		},
		"whitelist, input, but no filters: is not allowed": {
			inputs:         []string{"/projects/test"},
			filters:        nil,
			filterMode:     filterModeWhiteList,
			expectedResult: false,
		},
		"whitelist, input, but non matching filters: is not allowed": {
			inputs:         []string{"/projects/test"},
			filters:        []string{"repositories", "testrepo"},
			filterMode:     filterModeWhiteList,
			expectedResult: false,
		},
		"whitelist, filter exactly matches input: is allowed": {
			inputs:         []string{"/projects/test"},
			filters:        []string{"/projects/test", "non-matching-filter"},
			filterMode:     filterModeWhiteList,
			expectedResult: true,
		},
		"whitelist, regex filter matches input: is allowed": {
			inputs:         []string{"/projects/test"},
			filters:        []string{"non-matching-filter", "^.*test$"},
			filterMode:     filterModeWhiteList,
			expectedResult: true,
		},
		"whitelist, simple string as regex filter matches input: is allowed": {
			inputs:         []string{"/projects/test"},
			filters:        []string{"non-matching-filter", "test"},
			filterMode:     filterModeWhiteList,
			expectedResult: true,
		},
	}
	for testCaseName, testData := range testDataSet {
		t.Run(testCaseName, func(t *testing.T) {
			for inputIndex, input := range testData.inputs {
				f := newFilter(testData.filters, testData.filterMode)
				result := f.isAllowed(input)

				assert.Exactly(t, testData.expectedResult, result, fmt.Sprintf("inputIndex: %v", inputIndex))
			}
		})
	}
}
