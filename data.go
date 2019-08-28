package main

import (
	"strings"

	syscall "golang.org/x/sys/unix"
)

type erow struct {
	chars  []rune
	render []rune
}

type editorConfig struct {
	cx, cy      int
	rx          int
	screenRows  int
	screenCols  int
	rows        []erow
	rowOffset   int
	colOffset   int
	origTermios syscall.Termios
}

var cfg editorConfig

const kiloTabStop = 4

var tabSpaces = []rune(strings.Repeat(" ", kiloTabStop))
