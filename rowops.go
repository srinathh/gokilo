package main

type erow []rune

func (row erow) text() []rune {
	dest := []rune{}
	for _, r := range row {
		switch r {
		case '\t':
			dest = append(dest, tabSpaces...)
		default:
			dest = append(dest, r)
		}
	}
	return dest
}

// cxToRx transforms cursor positions to account for tab stops
func (row erow) cxToRx(cx int) int {
	rx := 0
	for j := 0; j < cx; j++ {
		if row[j] == '\t' {
			rx = (rx + kiloTabStop - 1) - (rx % kiloTabStop)
		}
		rx++
	}
	return rx

}

func (row erow) delChar(at int) []rune {
	if at < 0 || at >= len(row) {
		return row
	}

	copy(row[at:], row[at+1:])
	row = row[:len(row)-1]
	return row
}

// Insert Operations
func (row erow) insertChar(at, c int) []rune {
	// if at out of bounds, append to the end of the row
	if at < 0 || at > len(row) {
		return row
	}

	row = append(row, 0)
	copy(row[at+1:], row[at:])
	row[at] = rune(c)
	return row
}
