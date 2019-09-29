package main

import (
	"gokilo/internal/runes"
)

func editorFindCallback(query string, key int) {
	switch key {
	case '\r', '\x1b':
		cfg.lastMatch = -1
		cfg.direction = 1
		return

	case keyArrowRight, keyArrowDown:
		cfg.direction = 1
	case keyArrowLeft, keyArrowUp:
		cfg.direction = -1
	default:
		cfg.lastMatch = -1
		cfg.direction = 1
	}

	if cfg.lastMatch == -1 {
		cfg.direction = 1
	}

	current := cfg.lastMatch

	for i := 0; i < len(cfg.rows); i++ {

		current = current + cfg.direction
		if current == -1 {
			current = len(cfg.rows) - 1
		} else if current == len(cfg.rows) {
			current = 0
		}

		row := cfg.rows[current]

		if idx := runes.Index(runes.ToLower(row.chars), runes.ToLower([]rune(query))); idx != -1 {
			cfg.lastMatch = current
			cfg.cy = current
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

	query := editorPrompt("Search (use ESC/ARROWS/ENTER): %s", editorFindCallback)

	if query == "" {
		cfg.cx = savedCx
		cfg.cy = savedCy
		cfg.colOffset = savedColOffset
		cfg.rowOffset = savedRowOffset
	}

}
