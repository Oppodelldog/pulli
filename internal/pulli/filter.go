package pulli

import (
	"regexp"
)

const FilterModeWhiteList = "whitelist"
const FilterModeBlackList = "blacklist"

func newFilter(filters []string, filterMode string) *filter {
	if filterMode == "" {
		filterMode = FilterModeBlackList
	}

	regExFilters := buildRegExMatchers(filters)

	return &filter{
		filters:      filters,
		filterMode:   filterMode,
		regExFilters: regExFilters,
	}
}

type filter struct {
	filters      []string
	filterMode   string
	regExFilters []*regexp.Regexp
}

func (f *filter) isAllowed(path string) bool {
	isPathMatching := f.isPathMatchingFilter(path)

	return (isPathMatching && f.filterMode == FilterModeWhiteList) ||
		(!isPathMatching && f.filterMode == FilterModeBlackList)
}

func (f *filter) isPathMatchingFilter(path string) bool {
	return f.isPathMatchingString(path) || f.isPathMatchingRegExMatcher(path)
}

func (f *filter) isPathMatchingRegExMatcher(path string) bool {
	for _, matcher := range f.regExFilters {
		if matcher.MatchString(path) {
			return true
		}
	}

	return false
}

func (f *filter) isPathMatchingString(path string) bool {
	for _, exclude := range f.filters {
		if exclude == path {
			return true
		}
	}

	return false
}

func buildRegExMatchers(patterns []string) []*regexp.Regexp {
	var regExMatcher []*regexp.Regexp //nolint:prealloc
	for _, pattern := range patterns {
		regExMatcher = append(regExMatcher, regexp.MustCompile(pattern))
	}

	return regExMatcher
}
