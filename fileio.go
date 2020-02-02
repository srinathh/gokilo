package main

import (
	"bufio"
	"fmt"
	"os"
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

// Save writes a file to disk
func Save(rows []ERow, filename string) error {

	fil, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %s: %w", filename, err)
	}
	defer fil.Close()

	for _, row := range rows {
		if _, err := fmt.Fprintf(fil, "%s\n", string(row)); err != nil {
			return fmt.Errorf("error writing to file %s: %w", filename, err)
		}
	}

	if err = fil.Close(); err != nil {
		return fmt.Errorf("error closing written file: %s: %w", filename, err)
	}
	return nil
}
