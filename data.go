package main

import syscall "golang.org/x/sys/unix"

type erow []rune

type editorConfig struct {
	cx, cy      int
	screenRows  int
	screenCols  int
	rows        []erow
	origTermios syscall.Termios
}

var cfg editorConfig