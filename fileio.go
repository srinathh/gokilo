package main

import (
	"bufio"
	"os"
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
	cfg.fileName = fileName
	return nil
}

func editorUpdateRow(src []rune) []rune {
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
