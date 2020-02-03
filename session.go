package main

const kiloQuitTimes = 1

// Session stores state pertaining to the editing session
type Session struct {
	QuitTimes int // how many times Ctrl + Q has been pressed
	LastMatch int // find state
	Direction int // find direction
}

// NewSession initializes a new session
func NewSession(filename string) *Session {

	return &Session{
		QuitTimes: kiloQuitTimes,
		LastMatch: -1,
		Direction: 1,
	}
}
