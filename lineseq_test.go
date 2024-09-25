package seq

import (
	"errors"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLineSeqNext(t *testing.T) {
	f, err := os.Open("./petnames.txt")
	assert.Nil(t, err)
	lineSeq := NewLineSeq(f)
	line, err := lineSeq.Next()
	assert.Nil(t, err)
	assert.Equal(t, "AJ", line)
}

func TestLineSeqNext5(t *testing.T) {
	f, err := os.Open("./petnames.txt")
	assert.Nil(t, err)
	lineSeq := NewLineSeq(f)
	var line string
	for range 5 {
		line, err = lineSeq.Next()
		assert.Nil(t, err)
	}
	assert.Equal(t, "Abigail", line)
}

func TestLineSeqNext13(t *testing.T) {
	f, err := os.Open("./petnames.txt")
	assert.Nil(t, err)
	lineSeq := NewLineSeq(f)
	var line string
	for range 13 {
		line, err = lineSeq.Next()
		assert.Nil(t, err)
	}
	assert.Equal(t, "Alf", line)
}

func TestLineSeqThorough(t *testing.T) {
	str := "one\ntwo\nthree\n"
	rd := strings.NewReader(str)
	sq := NewLineSeq(rd)
	var line string
	var err error
	// one
	line, err = sq.Next()
	assert.Nil(t, err)
	assert.Nil(t, sq.Err())
	assert.Equal(t, "one", line)
	assert.Equal(t, 0, sq.LastPosition())
	assert.Equal(t, 4, sq.Position())
	// two
	line, err = sq.Next()
	assert.Nil(t, err)
	assert.Nil(t, sq.Err())
	assert.Equal(t, "two", line)
	assert.Equal(t, 4, sq.LastPosition())
	assert.Equal(t, 8, sq.Position())
	// three
	line, err = sq.Next()
	assert.Nil(t, err)
	assert.Nil(t, sq.Err())
	assert.Equal(t, "three", line)
	assert.Equal(t, 8, sq.LastPosition())
	assert.Equal(t, 14, sq.Position())
	// EOF
	line, err = sq.Next()
	assert.True(t, errors.Is(err, io.EOF))
	assert.True(t, errors.Is(sq.Err(), io.EOF))
	assert.Equal(t, "", line)
	assert.Equal(t, 14, sq.LastPosition())
	assert.Equal(t, 14, sq.Position())
	// Confirm EOF
	testEof(t, sq)
}

func TestLineSeqCount(t *testing.T) {
	f, err := os.Open("./petnames.txt")
	assert.Nil(t, err)
	sq := NewLineSeq(f)
	n, err := Count(sq)
	assert.Nil(t, err)
	assert.Equal(t, 1000, n)
	testEof(t, sq)
}
