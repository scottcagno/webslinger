package tests

func tolower(c byte) byte {
	if 'A' <= c && c <= 'Z' {
		c += 'a' - 'A'
	}
	return c
}

func MatchV3(pat string, s string) (bool, error) {
	return match(pat, s), nil
}

func match(pattern, name string) bool {
	px := 0
	nx := 0
	nextPx := 0
	nextNx := 0
	for px < len(pattern) || nx < len(name) {
		if px < len(pattern) {
			c := pattern[px]
			switch c {
			default: // ordinary character
				if nx < len(name) && name[nx] == c {
					px++
					nx++
					continue
				}
			case '?': // single-character wildcard
				if nx < len(name) {
					px++
					nx++
					continue
				}
			case '*': // zero-or-more-character wildcard
				// Try to match at nx.
				// If that doesn't work out,
				// restart at nx+1 next.
				nextPx = px
				nextNx = nx + 1
				px++
				continue
			}
		}
		// Mismatch. Maybe restart.
		if 0 < nextNx && nextNx <= len(name) {
			px = nextPx
			nx = nextNx
			continue
		}
		return false
	}
	// Matched all of pattern to all of name. Success.
	return true
}

func stringMatch(pattern []byte, s []byte, nocase bool) bool {
	patternLen := len(pattern)
	stringLen := len(s)
	for patternLen > 0 && stringLen > 0 {
		switch pattern[0] {
		case '*':
			for patternLen > 1 && pattern[1] == '*' {
				pattern = pattern[1:]
				patternLen--
			}
			if patternLen == 1 {
				return true
			}
			for stringLen > 0 {
				if stringMatch(pattern[1:], s, nocase) {
					return true
				}
				s = s[1:]
				stringLen--
			}
			return false
		case '?':
			s = s[1:]
			stringLen--
		case '[':
			var not, match bool

			pattern = pattern[1:]
			patternLen--
			not = pattern[0] == '^'
			if not {
				pattern = pattern[1:]
				patternLen--
			}
			match = false
			for {
				if pattern[0] == '\\' && patternLen >= 2 {
					pattern = pattern[1:]
					patternLen--
					if pattern[0] == s[0] {
						match = true
					}
				} else if pattern[0] == ']' {
					break
				} else if patternLen == 0 {
					pattern = pattern[len(pattern)-1:]
					patternLen++
					break
				} else if patternLen >= 3 && pattern[1] == '-' {
					start := pattern[0]
					end := pattern[2]
					c := s[0]
					if start > end {
						start, end = end, start
					}
					if nocase {
						start = tolower(start)
						end = tolower(end)
						c = tolower(c)
					}
					pattern = pattern[2:]
					patternLen -= 2
					if c >= start && c <= end {
						match = true
					}
				} else {
					if !nocase {
						if pattern[0] == s[0] {
							match = true
						}
					} else {
						if tolower(pattern[0]) == tolower(s[0]) {
							match = true
						}
					}
				}
				pattern = pattern[1:]
				patternLen--
			}
			if not {
				match = !match
			}
			if !match {
				return false
			}
			s = s[1:]
			stringLen--
		case '\\':
			if patternLen >= 2 {
				pattern = pattern[1:]
				patternLen--
			}
			fallthrough
		default:
			if !nocase {
				if pattern[0] != s[0] {
					return false
				}
			} else {
				if tolower(pattern[0]) != tolower(s[0]) {
					return false
				}
			}
			s = s[1:]
			stringLen--
			break
		}
		pattern = pattern[1:]
		patternLen--
		if stringLen == 0 {
			var i int
			for pattern[i] == '*' {
				pattern = pattern[1:]
				i++
				patternLen--
			}
			break
		}
	}
	return patternLen == 0 && stringLen == 0
}
