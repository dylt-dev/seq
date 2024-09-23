package seq

import (
	"errors"
	"io"
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArrayFlcCount(t *testing.T) {
	lines := []string{"a", "c", "b"}
	flc := NewArrayFiniteLineCollection(lines)
	n, err := flc.Count()
	assert.Nil(t, err)
	assert.Equal(t, len(lines), n)
}

func TestArrayFlcGetLine0(t *testing.T) {
	lines := []string{"a", "c", "b"}
	flc := NewArrayFiniteLineCollection(lines)
	i := 0
	line, err := flc.GetLine(i)
	assert.Nil(t, err)
	assert.Equal(t, lines[i], line)
}

func TestArrayFlcGetLine1(t *testing.T) {
	lines := []string{"a", "c", "b"}
	flc := NewArrayFiniteLineCollection(lines)
	i := len(lines) - 1
	line, err := flc.GetLine(i)
	assert.Nil(t, err)
	assert.Equal(t, lines[i], line)
}

func TestArrayFlcGetLine2(t *testing.T) {
	lines := []string{"a", "c", "b"}
	flc := NewArrayFiniteLineCollection(lines)
	i := rand.IntN(len(lines))
	line, err := flc.GetLine(i)
	assert.Nil(t, err)
	assert.Equal(t, lines[i], line)
}

func TestArrayFlcGetRandomLines(t *testing.T) {
	lines := []string{"a", "c", "b"}
	var flc FiniteLineCollection = NewArrayFiniteLineCollection(lines)
	for range 1000 {
		line, err := GetRandomLine(flc)
		assert.Nil(t, err)
		assert.Contains(t, lines, line)
	}
}

func TestFileFlcCount(t *testing.T) {
	path := "./petnames.txt"
	fileFlc := NewFileFlc(path)
	n, err := fileFlc.Count()
	assert.Nil(t, err)
	assert.Equal(t, 1000, n)
}

func TestFileFlcCountTwice(t *testing.T) {
	path := "./petnames.txt"
	fileFlc := NewFileFlc(path)
	n, err := fileFlc.Count()
	assert.Nil(t, err)
	assert.Equal(t, 1000, n)
	n, err = fileFlc.Count()
	assert.Nil(t, err)
	assert.Equal(t, 1000, n)
}

func TestFileFlcGetLine0(t *testing.T) {
	path := "./petnames.txt"
	fileFlc := NewFileFlc(path)
	i := 0
	target := "AJ"
	line, err := fileFlc.GetLine(i)
	assert.Nil(t, err)
	assert.Equal(t, target, line)
}

func TestFileFlcGetLine1(t *testing.T) {
	path := "./petnames.txt"
	fileFlc := NewFileFlc(path)
	i := 12
	target := "Alf"
	line, err := fileFlc.GetLine(i)
	assert.Nil(t, err)
	assert.Equal(t, target, line)
}

func TestFileFlcGetLine2(t *testing.T) {
	path := "./petnames.txt"
	fileFlc := NewFileFlc(path)
	i := 999
	target := "Zorro"
	line, err := fileFlc.GetLine(i)
	assert.Nil(t, err)
	assert.Equal(t, target, line)
}

func TestFileFlcGetRandomLines(t *testing.T) {
	path := "./petnames.txt"
	var fileFlc FiniteLineCollection = NewFileFlc(path)
	n := 1000
	for range n {
		_, err := GetRandomLine(fileFlc)
		assert.Nil(t, err)
	}
}

func TestRandomLineSeq0(t *testing.T) {
	lines := []string{"a", "b", "c"}
	var flc FiniteLineCollection = NewArrayFiniteLineCollection(lines)
	seq := NewRandomLineSeq(flc, 0)
	for {
		line, err := seq.Next()
		if line != "" || err == nil {
			t.Logf("line=%s\n", line)
		}
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			} else {
				t.Fatal(err.Error())
			}
		}
	}
}

func TestRandomLineSeq1(t *testing.T) {
	lines := []string{"a", "b", "c"}
	var flc FiniteLineCollection = NewArrayFiniteLineCollection(lines)
	seq := NewRandomLineSeq(flc, 1)
	for {
		line, err := seq.Next()
		if line != "" || err == nil {
			t.Logf("line=%s\n", line)
		}
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			} else {
				t.Fatal(err)
			}
		}
	}
}

func TestRandomLineSeq2(t *testing.T) {
	var flc FiniteLineCollection = NewFileFlc("./petnames.txt")
	seq := NewRandomLineSeq(flc, 0)
	for range 10000 {
		line, err := seq.Next()
		if line != "" || err == nil {
			t.Logf("%s\n", line)
		}
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			} else {
				t.Fatal(err)
			}
		}
	}
}

func TestRandomLineSeq3(t *testing.T) {
	lines := []string{"a", "b", "c"}
	var flc FiniteLineCollection = NewArrayFiniteLineCollection(lines)
	seq := NewRandomLineSeq(flc, 0)
	for i, line := range IterWithIndex(seq) {
		if seq.Err() == nil && line != "" {
			t.Logf("%d: line=%s\n", i, line)
		}
	}
}
