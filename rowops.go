package main

func editorUpdateRow(src []rune) []rune {
	dest := []rune{}
	for _, r := range src {
		switch r {
		case '\t':
			dest = append(dest, tabSpaces...)
		default:
			dest = append(dest, r)
		}
	}
	return dest
}

// editorRowCxToRx transforms cursor positions to account for tab stops
func editorRowCxToRx(row []rune, cx int) int {
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
		cfg.rows[cfg.cy].chars = editorRowDelChar(cfg.rows[cfg.cy].chars, cfg.cx-1)
		cfg.rows[cfg.cy].render = editorUpdateRow(cfg.rows[cfg.cy].chars)
		cfg.cx--
	} else {
		cfg.cx = len(cfg.rows[cfg.cy-1].chars)
		cfg.rows[cfg.cy-1].chars = append(cfg.rows[cfg.cy-1].chars, cfg.rows[cfg.cy].chars...)
		cfg.rows[cfg.cy-1].render = editorUpdateRow(cfg.rows[cfg.cy-1].chars)
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

/*
func editorRowDelChar(rowidx, at int) {
	if at < 0 || at >= len(cfg.rows[rowidx].chars) {
		return
	}

	copy(cfg.rows[rowidx].chars[at:], cfg.rows[rowidx].chars[at+1:])
	cfg.rows[rowidx].chars = cfg.rows[rowidx].chars[:len(cfg.rows[rowidx].chars)-1]

	cfg.rows[rowidx].render = editorUpdateRow(cfg.rows[rowidx].chars)
	cfg.dirty = true
}
*/

func editorDelRow(rowidx int) {
	if rowidx < 0 || rowidx >= len(cfg.rows) {
		return
	}

	copy(cfg.rows[rowidx:], cfg.rows[rowidx+1:])
	cfg.rows = cfg.rows[:len(cfg.rows)-1]
	cfg.dirty = true
}

/*
func editorRowAppendString(rowidx int, s []rune) {
	cfg.rows[rowidx].chars = append(cfg.rows[rowidx].chars, s...)
	cfg.dirty = true
	cfg.rows[rowidx].render = editorUpdateRow(cfg.rows[rowidx].chars)
}
*/

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

/*
func editorRowInsertChar(rowidx, at, c int) {
	// if at out of bounds, append to the end of the row
	if at < 0 || at > len(cfg.rows[rowidx].chars) {
		return
	}

	tmp := append(cfg.rows[rowidx].chars[:at], rune(c))
	if at < len(cfg.rows[rowidx].chars) {
		tmp = append(tmp, cfg.rows[rowidx].chars[at:]...)
	}
	cfg.rows[rowidx].chars = tmp

	editorUpdateRow(cfg.rows[rowidx].chars)
	cfg.dirty = true
}
*/

func editorInsertChar(c int) {
	if cfg.cy == len(cfg.rows) {
		editorInsertRow(len(cfg.rows), "")
	}
	cfg.rows[cfg.cy].chars = editorRowInsertChar(cfg.rows[cfg.cy].chars, cfg.cx, c)
	cfg.rows[cfg.cy].render = editorUpdateRow(cfg.rows[cfg.cy].chars)
	cfg.dirty = true
	cfg.cx++
}

/*
// adding and updating rows
func editorAppendRow(s string) {
	rns := []rune(s)
	row := erow{
		chars:  rns,
		render: editorUpdateRow(rns),
	}
	cfg.rows = append(cfg.rows, row)
	cfg.dirty = true
}
*/

func editorInsertNewline() {
	if cfg.cx == 0 {
		editorInsertRow(cfg.cy, "")
		return
	}

	moveChars := string(cfg.rows[cfg.cy].chars[cfg.cx:])

	cfg.rows[cfg.cy].chars = cfg.rows[cfg.cy].chars[:cfg.cx]
	cfg.rows[cfg.cy].render = editorUpdateRow(cfg.rows[cfg.cy].chars)

	editorInsertRow(cfg.cy+1, moveChars)

	cfg.cy++
	cfg.cx = 0
}

func editorInsertRow(rowidx int, s string) {
	if rowidx < 0 || rowidx > len(cfg.rows) {
		return
	}

	rns := []rune(s)
	row := erow{
		chars:  rns,
		render: editorUpdateRow(rns),
	}

	cfg.rows = append(cfg.rows, erow{})
	copy(cfg.rows[rowidx+1:], cfg.rows[rowidx:])
	cfg.rows[rowidx] = row

	cfg.dirty = true
}
