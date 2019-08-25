package main

import syscall "golang.org/x/sys/unix"

type erow struct {
	chars  []rune
	render []rune
}

type editorConfig struct {
	cx, cy      int
	screenRows  int
	screenCols  int
	rows        []erow
	rowOffset   int
	colOffset   int
	origTermios syscall.Termios
}

var cfg editorConfig
