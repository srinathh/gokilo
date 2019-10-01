package main

import (
	"gokilo/internal/runes"
)

func editorFindCallback(query string, key int) {
	switch key {
	case '\r', '\x1b':
		editor.LastMatch = -1
		editor.Direction = 1
		return

	case keyArrowRight, keyArrowDown:
		editor.Direction = 1
	case keyArrowLeft, keyArrowUp:
		editor.Direction = -1
	default:
		editor.LastMatch = -1
		editor.Direction = 1
	}

	if editor.LastMatch == -1 {
		editor.Direction = 1
	}

	current := editor.LastMatch

	for i := 0; i < len(editor.Rows); i++ {

		current = current + editor.Direction
		if current == -1 {
			current = len(editor.Rows) - 1
		} else if current == len(editor.Rows) {
			current = 0
		}

		row := editor.Rows[current]

		if idx := runes.Index(runes.ToLower(row), runes.ToLower([]rune(query))); idx != -1 {
			editor.LastMatch = current
			editor.Cy = current
			editor.Cx = idx
			editor.RowOffset = len(editor.Rows)
			break
		}

	}
}

func editorFind() {

	savedCx := editor.Cx
	savedCy := editor.Cy
	savedColOffset := editor.ColOffset
	savedRowOffset := editor.RowOffset

	query := editorPrompt("Search (use ESC/ARROWS/ENTER): %s", editorFindCallback)

	if query == "" {
		editor.Cx = savedCx
		editor.Cy = savedCy
		editor.ColOffset = savedColOffset
		editor.RowOffset = savedRowOffset
	}

}
