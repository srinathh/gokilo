package main

// editorRowCxToRx transforms cursor positions to account for tab stops
func editorRowCxToRx(row erow, cx int) int {
	rx := 0
	for j := 0; j < cx; j++ {
		if row[j] == '\t' {
			rx = (rx + kiloTabStop - 1) - (rx % kiloTabStop)
		}
		rx++
	}
	return rx

}

// Delete operations

func editorDelChar() {
	// if the cursor is in the empty line at the end, do nothing (why?)
	if cfg.cy == len(cfg.rows) {
		return
	}

	// if at the beginning of the text, then do nothing
	if cfg.cx == 0 && cfg.cy == 0 {
		return
	}

	// different handling for at the beginning of the line or middle of line
	if cfg.cx > 0 {
		cfg.rows[cfg.cy] = editorRowDelChar(cfg.rows[cfg.cy], cfg.cx-1)
		cfg.cx--
	} else {
		cfg.cx = len(cfg.rows[cfg.cy-1])
		cfg.rows[cfg.cy-1] = append(cfg.rows[cfg.cy-1], cfg.rows[cfg.cy]...)
		editorDelRow(cfg.cy)
		cfg.cy--
	}
	cfg.dirty = true
}

func editorRowDelChar(row []rune, at int) []rune {
	if at < 0 || at >= len(row) {
		return row
	}

	copy(row[at:], row[at+1:])
	row = row[:len(row)-1]
	return row
}

func editorDelRow(rowidx int) {
	if rowidx < 0 || rowidx >= len(cfg.rows) {
		return
	}

	copy(cfg.rows[rowidx:], cfg.rows[rowidx+1:])
	cfg.rows = cfg.rows[:len(cfg.rows)-1]
	cfg.dirty = true
}

// Insert Operations
func editorRowInsertChar(row []rune, at, c int) []rune {
	// if at out of bounds, append to the end of the row
	if at < 0 || at > len(row) {
		return row
	}

	row = append(row, 0)
	copy(row[at+1:], row[at:])
	row[at] = rune(c)
	return row
}

func editorInsertChar(c int) {
	if cfg.cy == len(cfg.rows) {
		editorInsertRow(len(cfg.rows), "")
	}
	cfg.rows[cfg.cy] = editorRowInsertChar(cfg.rows[cfg.cy], cfg.cx, c)
	cfg.dirty = true
	cfg.cx++
}

func editorInsertNewline() {
	if cfg.cx == 0 {
		editorInsertRow(cfg.cy, "")
		return
	}

	moveChars := string(cfg.rows[cfg.cy][cfg.cx:])

	cfg.rows[cfg.cy] = cfg.rows[cfg.cy][:cfg.cx]

	editorInsertRow(cfg.cy+1, moveChars)

	cfg.cy++
	cfg.cx = 0
}

func editorInsertRow(rowidx int, s string) {
	if rowidx < 0 || rowidx > len(cfg.rows) {
		return
	}

	row := []rune(s)

	cfg.rows = append(cfg.rows, erow{})
	copy(cfg.rows[rowidx+1:], cfg.rows[rowidx:])
	cfg.rows[rowidx] = row

	cfg.dirty = true
}
