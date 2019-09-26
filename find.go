package main

import (
	"gokilo/internal/runes"
)

func editorFindCallback(query string, key int) {
	if key == '\r' || key == '\x1b' {
		return
	}

	for i, row := range cfg.rows {

		if idx := runes.Index(row.chars, []rune(query)); idx != -1 {
			cfg.cy = i
			cfg.cx = idx
			cfg.rowOffset = len(cfg.rows)
			break
		}

	}
}

func editorFind() {

	savedCx := cfg.cx
	savedCy := cfg.cy
	savedColOffset := cfg.colOffset
	savedRowOffset := cfg.rowOffset

	query := editorPrompt("Search (ESC to cancel): %s", editorFindCallback)

	if query == "" {
		cfg.cx = savedCx
		cfg.cy = savedCy
		cfg.colOffset = savedColOffset
		cfg.rowOffset = savedRowOffset
	}

}
