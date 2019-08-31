package main

import (
	"strings"
	"time"

	syscall "golang.org/x/sys/unix"
)

type erow struct {
	chars  []rune
	render []rune
}

type editorConfig struct {
	cx, cy        int
	rx            int
	screenRows    int
	screenCols    int
	rows          []erow
	rowOffset     int
	colOffset     int
	origTermios   syscall.Termios
	fileName      string
	statusMsg     string
	statusMsgTime time.Time
}

var cfg editorConfig

const kiloTabStop = 4

var tabSpaces = []rune(strings.Repeat(" ", kiloTabStop))
