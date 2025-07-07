package goprompt

// FilterRanked combines prefix, contains, and fuzzy filters in that order to create suggestions that
// consider the most relevant matches first.
func filterCombined(completions []Suggest, sub string, ignoreCase bool) []Suggest {
	prefixMatches := FilterHasPrefix(completions, sub, ignoreCase)
	containsMatches := FilterContains(completions, sub, ignoreCase)
	fuzzyMatches := FilterFuzzy(completions, sub, ignoreCase)

	// combine all matches, ensuring no duplicates
	presenseSet := make(map[string]Suggest)

	res := make([]Suggest, 0, len(prefixMatches)+len(containsMatches)+len(fuzzyMatches))

	for _, match := range prefixMatches {
		presenseSet[match.Text] = match
		res = append(res, match)
	}

	for _, match := range containsMatches {
		if _, exists := presenseSet[match.Text]; !exists {
			presenseSet[match.Text] = match
			res = append(res, match)
		}
	}

	for _, match := range fuzzyMatches {
		if _, exists := presenseSet[match.Text]; !exists {
			res = append(res, match)
		}
	}

	return res
}
