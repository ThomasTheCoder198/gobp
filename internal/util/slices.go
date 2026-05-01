package util

import "slices"

// StringIn reports whether target appears in ss.
func StringIn(ss []string, target string) bool {
	return slices.Contains(ss, target)
}
