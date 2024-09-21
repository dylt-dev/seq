package seq

import (
	"errors"
	"io"
	"math/rand/v2"
	"os"
)

type FiniteLineCollection interface {
	Count() (int, error)
	GetLine(int) (string, error)
}

func GetRandomLine(flc FiniteLineCollection) (string, error) {
	n, err := flc.Count()
	if err != nil {
		return "", err
	}
	i := rand.IntN(n)
	line, err := flc.GetLine(i)
	return line, err
}

type ArrayFiniteLineCollection struct {
	lines []string
}

func NewArrayFiniteLineCollection(lines []string) *ArrayFiniteLineCollection {
	return &ArrayFiniteLineCollection{
		lines,
	}
}

func (flc *ArrayFiniteLineCollection) Count() (int, error) {
	return len(flc.lines), nil
}

func (flc *ArrayFiniteLineCollection) GetLine(i int) (string, error) {
	return flc.lines[i], nil
}

type FileFlc struct {
	path string
}

func NewFileFlc(path string) *FileFlc {
	return &FileFlc{path}
}

func (flc *FileFlc) Count() (int, error) {
	f, err := os.Open(flc.path)
	if err != nil {
		return 0, err
	}
	var seq Seq[string] = NewLineSeq(f)
	n, err := Count(seq)
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return n, err
}

func (flc *FileFlc) GetLine(n int) (string, error) {
	f, err := os.Open(flc.path)
	if err != nil {
		return "", err
	}
	seq := NewLineSeq(f)
	var line string
	for range n + 1 {
		line, err = seq.Next()
		if err != nil && line == "" {
			return "", err
		}
	}
	return line, nil
}
