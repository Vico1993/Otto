package utils

import "strings"

// Find a string inside a slice if present
func InSlice(str string, slice []string) bool {
	str = strings.ToLower(str)

	for _, s := range slice {
		s := strings.ToLower(s)

		if s == str || strings.Contains(s, str) {
			return true
		}
	}
	return false
}
