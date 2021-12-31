package httpsfv

// isLowerCaseAlpha checks if c is a lower cased alpha character.
func isLowerCaseAlpha(c byte) bool {
	return 'a' <= c && c <= 'z'
}

// isAlpha checks if c is an alpha character.
func isAlpha(c byte) bool {
	return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')
}

// isDigit checks if c is a digit.
func isDigit(c byte) bool {
	return '0' <= c && c <= '9'
}
