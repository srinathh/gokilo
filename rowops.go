package main

func editorAppendRow(s string) {
	rns := []rune(s)
	row := erow{
		chars:  rns,
		render: editorUpdateRow(rns),
	}
	cfg.rows = append(cfg.rows, row)
	cfg.dirty = true
}

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

func editorRowDelChar(rowidx, at int) {
	if at < 0 || at >= len(cfg.rows[rowidx].chars) {
		return
	}
	copy(cfg.rows[rowidx].chars[at:], cfg.rows[rowidx].chars[at+1:])
	cfg.rows[rowidx].chars = cfg.rows[rowidx].chars[:len(cfg.rows[rowidx].chars)-1]
	cfg.rows[rowidx].render = editorUpdateRow(cfg.rows[rowidx].chars)
	cfg.dirty = true
}

//
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
