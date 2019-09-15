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
	cfg.rows[cfg.cy].chars = editorRowInsertChar(cfg.rows[cfg.cy].chars, cfg.cx, c)
	cfg.rows[cfg.cy].render = editorUpdateRow(cfg.rows[cfg.cy].chars)
	cfg.dirty = true
	cfg.cx++
}

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

func editorPrompt(prompt string) string {

	buf := ""

	for {
		editorSetStatusMsg(prompt, buf)
		editorRefreshScreen()

		c := editorReadKey()
		switch c {
		case keyDelete, keyBackSpace, ctrlKey('h'):
			if len(buf) > 0 {
				buf = buf[:len(buf)-1]
			}
		case '\x1b':
			editorSetStatusMsg("")
			return ""
		case '\r':
			if len(buf) != 0 {
				editorSetStatusMsg("")
				return buf
			}
		default:
			if c != ctrlKey('c') && c < 128 {
				buf = buf + string(c)
			}
		}
	}
}
