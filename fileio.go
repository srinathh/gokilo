package main

import (
	"bufio"
	"os"
	"strings"
)

func editorOpen(fileName string) error {

	r, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer r.Close()

	cfg.rows = []erow{}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		row := erow{
			chars: []rune(scanner.Text()),
		}
		row.render = editorUpdateRow(row.chars)

		cfg.rows = append(cfg.rows, row)
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func editorUpdateRow(src []rune) []rune {
	dest := []rune{}
	for _, r := range src {
		if r == '\t' {
			dest = append(dest, []rune(strings.Repeat(" ", 8))...)
		}
		dest = append(dest, r)
	}
	return dest
}
