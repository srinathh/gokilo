package main

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

/*
func editorInsertRow(rowidx int, s string) {
	if rowidx < 0 || rowidx > len(cfg.rows) {
		return
	}

	rns := []rune(s)
	row := erow{
		chars:  rns,
		render: editorUpdateRow(rns),
	}

	tmp := append(cfg.rows[:rowidx], row)
	if rowidx < len(cfg.rows) {
		tmp = append(tmp, cfg.rows[rowidx:]...)
	}
	cfg.rows = tmp

	cfg.dirty = true

}
*/

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
func editorRowCxToRx(rowIdx, cx int) int {
	rx := 0
	for j := 0; j < cx; j++ {
		if cfg.rows[rowIdx].chars[j] == '\t' {
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
		editorRowDelChar(cfg.cy, cfg.cx-1)
		cfg.cx--
	} else {
		cfg.cx = len(cfg.rows[cfg.cy-1].chars)
		editorRowAppendString(cfg.cy-1, cfg.rows[cfg.cy].chars)
		editorDelRow(cfg.cy)
		cfg.cy--
	}
}

func editorRowDelChar(rowidx, at int) {
	if at < 0 || at >= len(cfg.rows[rowidx].chars) {
		return
	}

	copy(cfg.rows[rowidx].chars[at:], cfg.rows[rowidx].chars[at+1:])
	cfg.rows[rowidx].chars = cfg.rows[rowidx].chars[:len(cfg.rows[rowidx].chars)-1]

	cfg.rows[rowidx].render = editorUpdateRow(cfg.rows[rowidx].chars)
	cfg.dirty = true
}

func editorDelRow(rowidx int) {
	if rowidx < 0 || rowidx >= len(cfg.rows) {
		return
	}

	copy(cfg.rows[rowidx:], cfg.rows[rowidx+1:])
	cfg.rows = cfg.rows[:len(cfg.rows)-1]
	cfg.dirty = true
}

func editorRowAppendString(rowidx int, s []rune) {
	cfg.rows[rowidx].chars = append(cfg.rows[rowidx].chars, s...)
	cfg.dirty = true
	cfg.rows[rowidx].render = editorUpdateRow(cfg.rows[rowidx].chars)
}

// Insert Operations

/*
func editorRowInsertChar(rowidx, at, c int) {

	// if at out of bounds, append to the end of the row
	if at < 0 || at > len(cfg.rows[rowidx].chars) {
		cfg.rows[rowidx].chars = append(cfg.rows[rowidx].chars, rune(c))
		editorUpdateRow(cfg.rows[rowidx].chars)
		return
	}

	// else insert without additonal allocation
	cfg.rows[rowidx].chars = append(cfg.rows[rowidx].chars, 0)
	copy(cfg.rows[rowidx].chars[at+1:], cfg.rows[rowidx].chars[at:])
	cfg.rows[rowidx].chars[at] = rune(c)
	editorUpdateRow(cfg.rows[rowidx].chars)
	cfg.dirty = true

}
*/

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

func editorInsertChar(c int) {
	if cfg.cy == len(cfg.rows) {
		cfg.rows = append(cfg.rows, newErow())
	}
	editorRowInsertChar(cfg.cy, cfg.cx, c)
	cfg.rows[cfg.cy].render = editorUpdateRow(cfg.rows[cfg.cy].chars)
	cfg.cx++
}
