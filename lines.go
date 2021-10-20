package fileutils

import (
	"bufio"
	"bytes"
	"io"
	"os"
)

const bufSize = 32 * 1024

// LinesCount returns count of lines in the file.
func LinesCount(path string) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	return linesCount(file)
}

// linesCount count lines in provided io.Reader.
func linesCount(r io.Reader) (int, error) {
	buf := make([]byte, bufSize)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		// If there is no "\n" on the last line.
		if c != 0 && c < bufSize && buf[c-1] != '\n' {
			count++
		}

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

// FirstLine returns first line of the file.
func FirstLine(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if line := scanner.Text(); line != "" {
			return line, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", err
}

// LinesIterator represents file lines iterator.
type LinesIterator struct {
	path       string
	linesCount int
	scanner    *bufio.Scanner
	file       *os.File
}

// NewLinesIterator returns new file reader.
func NewLinesIterator(path string) (*LinesIterator, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	linesCount, err := linesCount(file)
	if err != nil {
		return nil, err
	}

	if _, err = file.Seek(0, os.SEEK_SET); err != nil {
		return nil, err
	}

	return &LinesIterator{
		path:       path,
		linesCount: linesCount,
		file:       file,
		scanner:    bufio.NewScanner(file),
	}, nil
}

// Next returns next read line.
func (r *LinesIterator) Next(line *string) bool {

	for r.scanner.Scan() {
		*line = r.scanner.Text()
		return true
	}

	return false
}

// Close closes file reader.
func (r *LinesIterator) Close() error {
	return r.file.Close()
}

// Reset resets file reader.
func (r *LinesIterator) Reset() {
	r.file.Seek(0, os.SEEK_SET)
	r.scanner = bufio.NewScanner(r.file)
}

// Count returns lines count in the reader.
func (r *LinesIterator) Count() uint64 {
	return uint64(r.linesCount)
}
