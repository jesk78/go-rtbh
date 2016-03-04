package pipeline

import (
	"regexp"
)

func MatchesRule(value string, re *regexp.Regexp) bool {
	result := re.Match([]byte(value))
	return result
}

func FoundMatch(value string) bool {
	for _, re := range Ruleset {
		if MatchesRule(value, re) {
			return true
		}
	}

	return false
}
