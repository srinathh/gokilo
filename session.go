package main

import "fmt"

const kiloQuitTimes = 1

// Session stores state pertaining to the editing session
type Session struct {
	Dirty     bool   // has the file been edited
	FileName  string // the path to the file being edited. Could be empty string
	QuitTimes int    // how many times Ctrl + Q has been pressed
	LastMatch int    // find state
	Direction int    // find direction
	Editor    *Editor
}

// NewSession initializes a new session
func NewSession(filename string) *Session {

	rows := []ERow{}

	if filename != "" {
		var err error
		if rows, err = Open(filename); err != nil {
			SafeExit(fmt.Errorf("Error opening file %s: %v", filename, err))
		}
	}

	return &Session{
		Dirty:     false,
		FileName:  filename,
		QuitTimes: kiloQuitTimes,
		LastMatch: -1,
		Editor: &Editor{
			Cx:   0,
			Cy:   0,
			Rows: rows,
		},
	}
}
