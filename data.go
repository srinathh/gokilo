package main

import (
	"strings"
	"time"
)

type erow []rune

func (src erow) Text() []rune {
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

/*
type erow struct {
	chars  []rune
	render []rune
}

func newErow() erow {
	return erow{
		[]rune{},
		[]rune{},
	}
}
*/

type editorConfig struct {
	cx, cy     int
	rx         int
	screenRows int
	screenCols int
	rows       []erow
	rowOffset  int
	colOffset  int
	//origTermios   syscall.Termios
	origTermCfg   interface{}
	fileName      string
	statusMsg     string
	statusMsgTime time.Time
	dirty         bool
	quitTimes     int
	lastMatch     int
	direction     int
}

var cfg editorConfig

const kiloTabStop = 4
const kiloQuitTimes = 1

var tabSpaces = []rune(strings.Repeat(" ", kiloTabStop))
