package main

import (
	"strings"
)

// Config collects configuraiton params from the app
type Config struct {
	ScreenRows  int
	ScreenCols  int
	OrigTermCfg []byte
}

const kiloTabStop = 4
const kiloQuitTimes = 1

var tabSpaces = []rune(strings.Repeat(" ", kiloTabStop))
