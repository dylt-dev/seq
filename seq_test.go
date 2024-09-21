package seq

import (
	_ "embed"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

//go:embed iter-test-0.txt
var TS0 string

func TestRuneSeq0(t *testing.T) {
	str := TS0
	sr := strings.NewReader(str)
	seq := NewRuneSeq(sr)
	for _, ru := range Iter2(seq) {
		fmt.Printf("ru=%c\n", ru)
	}
}

func TestRuneSeqCount(t *testing.T) {
	str := TS0
	sr := strings.NewReader(str)
	seq := NewRuneSeq(sr)
	n, err := Count(seq)
	assert.Equal(t, len(str), n)
	assert.Nil(t, err)
}

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
	seq := NewLineSeq(rd)
	var line string
	var err error
	// one
	line, err = seq.Next()
	assert.Nil(t, err)
	assert.Nil(t, seq.Err())
	assert.Equal(t, "one", line)
	assert.Equal(t, 0, seq.LastPosition())
	assert.Equal(t, 4, seq.Position())
	// two
	line, err = seq.Next()
	assert.Nil(t, err)
	assert.Nil(t, seq.Err())
	assert.Equal(t, "two", line)
	assert.Equal(t, 4, seq.LastPosition())
	assert.Equal(t, 8, seq.Position())
	// three
	line, err = seq.Next()
	assert.Nil(t, err)
	assert.Nil(t, seq.Err())
	assert.Equal(t, "three", line)
	assert.Equal(t, 8, seq.LastPosition())
	assert.Equal(t, 14, seq.Position())
	// EOF
	line, err = seq.Next()
	assert.True(t, errors.Is(err, io.EOF))
	assert.True(t, errors.Is(seq.Err(), io.EOF))
	assert.Equal(t, "", line)
	assert.Equal(t, 14, seq.LastPosition())
	assert.Equal(t, 14, seq.Position())
}

func TestLineSeqCount(t *testing.T) {
	f, err := os.Open("./petnames.txt")
	assert.Nil(t, err)
	lineSeq := NewLineSeq(f)
	n, err := Count(lineSeq)
	assert.Nil(t, err)
	assert.Equal(t, 1000, n)

}

func TestSimpleRacer(t *testing.T) {
	speed := float32(2.27)
	racer := NewSimpleRacer(speed)
	for i, dx := range Iter2[float32](racer) {
		t.Logf("%f", dx)
		time.Sleep(100 * time.Millisecond)
		if i >= 100 { break }
	}
}
