package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Open reads a file and returns erows representing each line
func Open(fileName string) ([]ERow, error) {

	r, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("error opening file %s: %w", fileName, err)
	}
	defer r.Close()

	ret := []ERow{}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		ret = append(ret, []rune(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file %s: %w", fileName, err)
	}

	return ret, err

}

func editorRowsToString() string {
	var sb strings.Builder

	for _, rows := range editor.Rows {
		sb.WriteString(string(rows))
		sb.WriteByte('\n')
	}
	return sb.String()
}

func editorSave() {
	if editor.FileName == "" {
		editor.FileName = editorPrompt("Save as: %s", nil)
		if editor.FileName == "" {
			editorSetStatusMsg("Save aborted!")
			return
		}
	}

	fil, err := os.Create(editor.FileName)
	if err != nil {
		editorSetStatusMsg("ERROR creating file: %s: %s", err, editor.FileName)
		return
	}
	defer fil.Close()

	if _, err := fmt.Fprint(fil, editorRowsToString()); err != nil {
		editorSetStatusMsg("ERROR writing to file: %s: %s", err, editor.FileName)
		return
	}

	if err = fil.Close(); err != nil {
		editorSetStatusMsg("ERROR closing written file: %s: %s", err, editor.FileName)
		return
	}

	editorSetStatusMsg("SAVED to file: %s", editor.FileName)
	editor.Dirty = false

}
