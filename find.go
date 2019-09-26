package main

import (
	"gokilo/internal/runes"
)

func editorFind() {
	query := editorPrompt("Search (ESC to cancel): %s")

	if query == "" {
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
