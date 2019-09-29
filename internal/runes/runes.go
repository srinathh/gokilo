package runes

import "unicode"

// IndexRune returns the index of the first instance of the Unicode code point
// r, or -1 if rune is not present in s.
func IndexRune(s []rune, c rune) int {
	for i, b := range s {
		if b == c {
			return i
		}
	}
	return -1
}

// Equal reports whether a and b
// are the same length and contain the same bytes.
// A nil argument is equivalent to an empty slice.
func Equal(a, b []rune) bool {
	if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// Index returns the index of the first instance of substr in s, or -1 if substr is not present in s.
func Index(s, substr []rune) int {
	n := len(substr)
	switch {
	case n == 0:
		return 0
	case n == 1:
		return IndexRune(s, substr[0])
	case n == len(s):
		if Equal(substr, s) {
			return 0
		}
		return -1
	case n > len(s):
		return -1
	default:
		return indexRabinKarp(s, substr)
	}

}

func indexRabinKarp(s, sep []rune) int {
	// Rabin-Karp search
	hashsep, pow := hashStr(sep)
	n := len(sep)
	var h uint32
	for i := 0; i < n; i++ {
		h = h*primeRK + uint32(s[i])
	}
	if h == hashsep && Equal(s[:n], sep) {
		return 0
	}
	for i := n; i < len(s); {
		h *= primeRK
		h += uint32(s[i])
		h -= pow * uint32(s[i-n])
		i++
		if h == hashsep && Equal(s[i-n:i], sep) {
			return i - n
		}
	}
	return -1
}

// primeRK is the prime base used in Rabin-Karp algorithm.
const primeRK = 16777619

// hashStr returns the hash and the appropriate multiplicative
// factor for use in Rabin-Karp algorithm.
func hashStr(sep []rune) (uint32, uint32) {
	hash := uint32(0)
	for i := 0; i < len(sep); i++ {
		hash = hash*primeRK + uint32(sep[i])
	}
	var pow, sq uint32 = 1, primeRK
	for i := len(sep); i > 0; i >>= 1 {
		if i&1 != 0 {
			pow *= sq
		}
		sq *= sq
	}
	return hash, pow
}

// Map returns a copy of the rune slice s with all its characters modified
// according to the mapping function. If mapping returns a negative value, the character is
// dropped from the byte slice with no replacement. The characters in s and the
// output are interpreted as UTF-8-encoded code points.
func Map(mapping func(r rune) rune, s []rune) []rune {
	ret := make([]rune, len(s))
	for j, r := range s {
		ret[j] = mapping(r)
	}
	return ret
}

// ToLower returns a copy of the rune slice s with all Unicode letters mapped to their lower case.
func ToLower(s []rune) []rune { return Map(unicode.ToLower, s) }

// ToUpper returns a copy of the rune slice s with all Unicode letters mapped to their lower case.
func ToUpper(s []rune) []rune { return Map(unicode.ToUpper, s) }

// ToTitle treats s as UTF-8-encoded bytes and returns a copy with all the Unicode letters mapped to their title case.
func ToTitle(s []rune) []rune { return Map(unicode.ToTitle, s) }
