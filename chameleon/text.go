package chameleon

import (
	"regexp"
	"strings"
)

var rexSafeChars = regexp.MustCompile(`[^a-zA-Z\-\.\_\/ ]`)

// safeText sanitizes the input to be safe to pass into system calls
func safeText(s string) string {
	s = rexSafeChars.ReplaceAllString(s, " ")
	for strings.Contains(s, "  ") {
		s = strings.ReplaceAll(s, "  ", " ")
	}
	s = strings.Trim(s, " ")
	s = strings.Trim(s, "/")
	for strings.Contains(s, "..") {
		s = strings.ReplaceAll(s, "..", ".")
	}
	return s
}
