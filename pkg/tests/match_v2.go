package tests

import (
	"bytes"
	"regexp"
	"strings"
)

// MatchV2 reports whether name matches the shell pattern.
func MatchV2(pattern, name string) (bool, error) {
	// Preprocess the pattern to convert it into a regular expression
	pat, err := preprocessPatternv2(pattern)
	if err != nil {
		return false, err
	}
	// Use the compiled regular expression to match against the name
	re, err := regexp.Compile(pat)
	if err != nil {
		return false, err
	}
	if ok := re.MatchString(name); !ok {
		return false, nil
	}
	return true, nil
}

// preprocessPattern takes a shell pattern and returns a regular expression equivalent
func preprocessPattern(pattern string) (string, error) {
	var b strings.Builder
	b.WriteByte('^')

	i := 0
	for i < len(pattern) {
		switch pattern[i] {
		case '*':
			b.WriteString(".*")
		case '?':
			b.WriteString(".")
		case '[':
			j := i + 1
			if j >= len(pattern) {
				return "", ErrBadPattern
			}
			if pattern[j] == '^' {
				j++
			}
			for j < len(pattern) && pattern[j] != ']' {
				if pattern[j] == '\\' && j+1 < len(pattern) {
					j++
				}
				j++
			}
			if j >= len(pattern) {
				return "", ErrBadPattern
			}
			b.WriteByte('[')
			if pattern[i+1] == '^' {
				b.WriteByte('^')
				i++
			}
			for k := i + 1; k < j; k++ {
				if pattern[k] == '\\' && k+1 < j {
					k++
					b.WriteByte(pattern[k])
				} else {
					b.WriteByte(pattern[k])
				}
			}
			b.WriteByte(']')
			i = j
		case '\\':
			i++
			if i >= len(pattern) {
				return "", ErrBadPattern
			}
			b.WriteByte(pattern[i])
		default:
			b.WriteByte(pattern[i])
		}
		i++
	}

	b.WriteByte('$')
	return b.String(), nil
}

func preprocessPatternv2(pattern string) (string, error) {
	buf := bytes.NewBufferString("^")

	i := 0
	for i < len(pattern) {
		switch pattern[i] {
		case '*':
			buf.WriteString(".*")
		case '?':
			buf.WriteByte('.')
		case '[':
			j := i + 1
			if j >= len(pattern) {
				return "", ErrBadPattern
			}
			if pattern[j] == '^' {
				j++
			}
			for j < len(pattern) && pattern[j] != ']' {
				if pattern[j] == '\\' && j+1 < len(pattern) {
					j++
				}
				j++
			}
			if j >= len(pattern) {
				return "", ErrBadPattern
			}
			buf.WriteByte('[')
			if pattern[i+1] == '^' {
				buf.WriteByte('^')
				i++
			}
			for k := i + 1; k < j; k++ {
				if pattern[k] == '\\' && k+1 < j {
					k++
					buf.WriteByte(pattern[k])
				} else {
					buf.WriteByte(pattern[k])
				}
			}
			buf.WriteByte(']')
			i = j
		case '\\':
			i++
			if i >= len(pattern) {
				return "", ErrBadPattern
			}
			buf.WriteByte(pattern[i])
		default:
			buf.WriteByte(pattern[i])
		}
		i++
	}

	buf.WriteByte('$')
	return buf.String(), nil
}
