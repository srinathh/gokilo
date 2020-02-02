package main

// ERow represents a line of text in a file
type ERow []rune

// Editor represents the data in the being edited in memory
type Editor struct {
	Cx, Cy int    // Cx and Cy represent current cursor position
	Rows   []ERow // Rows represent the textual data
}
