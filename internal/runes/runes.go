package runes

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

func Dummy() string {
	return "hello"
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

/*
// Index returns the index of the first instance of substr in s, or -1 if substr is not present in s.
func Index(s, substr []rune) int {
	n := len(substr)
	switch {
	case n == 0:
		return 0
	case n == 1:
		return IndexRune(s, substr[0])
	case n == len(s):
		if substr == s {
			return 0
		}
		return -1
	case n > len(s):
		return -1
	case n <= bytealg.MaxLen:
		// Use brute force when s and substr both are small
		if len(s) <= bytealg.MaxBruteForce {
			return bytealg.IndexString(s, substr)
		}
		c0 := substr[0]
		c1 := substr[1]
		i := 0
		t := len(s) - n + 1
		fails := 0
		for i < t {
			if s[i] != c0 {
				// IndexByte is faster than bytealg.IndexString, so use it as long as
				// we're not getting lots of false positives.
				o := IndexByte(s[i:t], c0)
				if o < 0 {
					return -1
				}
				i += o
			}
			if s[i+1] == c1 && s[i:i+n] == substr {
				return i
			}
			fails++
			i++
			// Switch to bytealg.IndexString when IndexByte produces too many false positives.
			if fails > bytealg.Cutover(i) {
				r := bytealg.IndexString(s[i:], substr)
				if r >= 0 {
					return r + i
				}
				return -1
			}
		}
		return -1
	}
	c0 := substr[0]
	c1 := substr[1]
	i := 0
	t := len(s) - n + 1
	fails := 0
	for i < t {
		if s[i] != c0 {
			o := IndexByte(s[i:t], c0)
			if o < 0 {
				return -1
			}
			i += o
		}
		if s[i+1] == c1 && s[i:i+n] == substr {
			return i
		}
		i++
		fails++
		if fails >= 4+i>>4 && i < t {
			// See comment in ../bytes/bytes_generic.go.
			j := indexRabinKarp(s[i:], substr)
			if j < 0 {
				return -1
			}
			return i + j
		}
	}
	return -1
}
*/
