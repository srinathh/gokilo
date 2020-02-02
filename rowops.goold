package main

// ERow represents a line of text in a file
type ERow []rune

// Text expands tabs in an eRow to spaces
func (row ERow) Text() []rune {
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

func spaces(n int) []rune {
	ret := make([]rune, n)
	for j := 0; j < n; j++ {
		ret[j] = ' '
	}
	return ret
}

// ScreenText returns a line of text offset to match screen window
func (row ERow) ScreenText(colOffset, screenWidth int) []rune {

	ret := spaces(screenWidth)

	txt := row.Text()
	if colOffset > len(txt) {
		return spaces(screenWidth)
	}

	for j := colOffset; j < len(txt) && j < colOffset+screenWidth; j++ {
		ret[j-colOffset] = txt[j]
	}
	return ret

}

// CxToRx transforms cursor positions to account for tab stops
func (row ERow) CxToRx(cx int) int {
	rx := 0
	for j := 0; j < cx; j++ {
		if row[j] == '\t' {
			rx = (rx + kiloTabStop - 1) - (rx % kiloTabStop)
		}
		rx++
	}
	return rx

}

// DelChar deletes a char at a position in a row
func (row ERow) DelChar(at int) []rune {
	if at < 0 || at >= len(row) {
		return row
	}

	copy(row[at:], row[at+1:])
	row = row[:len(row)-1]
	return row
}

// InsertChar inserts a rune at a position in an eRow
func (row ERow) InsertChar(at int, c rune) []rune {
	// if at out of bounds, append to the end of the row
	if at < 0 || at > len(row) {
		return row
	}

	row = append(row, 0)
	copy(row[at+1:], row[at:])
	row[at] = rune(c)
	return row
}
