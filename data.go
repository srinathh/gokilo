package main

import (
	"strings"
	"time"
)

// Config collects configuraiton params frot he app
type Config struct {
	ScreenRows    int
	ScreenCols    int
	OrigTermCfg   interface{}
	StatusMsg     string
	StatusMsgTime time.Time
}

const kiloTabStop = 4
const kiloQuitTimes = 1

var tabSpaces = []rune(strings.Repeat(" ", kiloTabStop))
