package main

func editorInsertChar(c int) {
	if cfg.cy == len(cfg.rows) {
		cfg.rows = append(cfg.rows, newErow())
	}
	editorRowInsertChar(cfg.cy, cfg.cx, c)
	cfg.rows[cfg.cy].render = editorUpdateRow(cfg.rows[cfg.cy].chars)
	cfg.cx++
}

func editorDelChar() {
	if cfg.cy == len(cfg.rows) {
		return
	}

	if cfg.cx > 0 {
		editorRowDelChar(cfg.cy, cfg.cx)
		cfg.cx--
	}
}
