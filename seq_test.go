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
	for _, ru := range IterWithIndex(seq) {
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
	for i, dx := range IterWithIndex(racer) {
		t.Logf("%f", dx)
		time.Sleep(100 * time.Millisecond)
		if i >= 100 {
			break
		}
	}
}

func TestLimitNextThorough(t *testing.T) {
	var err error
	rd, err := os.Open("./petnames.txt")
	if err != nil {
		panic("Unable to open petnames file")
	}
	petNamesSeq := NewLineSeq(rd)
	// names := []string{}
	newSeq := Limit(petNamesSeq, 5)
	var name, expected string
	// 0: "AJ"
	expected = "AJ"
	name, err = newSeq.Next()
	assert.Nil(t, err)
	assert.NotEmpty(t, name)
	assert.Equal(t, expected, name)
	// 1: "Abbey"
	expected = "Abbey"
	name, err = newSeq.Next()
	assert.Nil(t, err)
	assert.NotEmpty(t, name)
	assert.Equal(t, expected, name)
	// 2: "Abbie"
	expected = "Abbie"
	name, err = newSeq.Next()
	assert.Nil(t, err)
	assert.NotEmpty(t, name)
	assert.Equal(t, expected, name)
	// 3: "Abel"
	expected = "Abel"
	name, err = newSeq.Next()
	assert.Nil(t, err)
	assert.NotEmpty(t, name)
	assert.Equal(t, expected, name)
	// 4: "Abigail"
	expected = "Abigail"
	name, err = newSeq.Next()
	assert.Nil(t, err)
	assert.NotEmpty(t, name)
	assert.Equal(t, expected, name)
	// 5: "", ni
	expected = ""
	name, err = newSeq.Next()
	assert.NotNil(t, err)
	assert.True(t, errors.Is(err, io.EOF))
	assert.Empty(t, name)

	// for name := range Iter1(Limit(petNamesSeq, 5)) {
	// 	// Sequences let you check if the previous read actually resulted in an error, like EOF
	// 	if petNamesSeq.Err() == nil {
	// 		fmt.Printf("name=%s\n", name)
	// 		names = append(names, name)
	// 	} else {
	// 		// EOFs are ok. Other errors are terrifying.
	// 		if errors.Is(petNamesSeq.Err(), io.EOF) {
	// 			fmt.Println("Normal EOF reached. All is well.")
	// 		} else {
	// 			fmt.Printf("%v\n", err)
	// 			panic("Unepxected error!")
	// 		}
	// 	}
	// }
	// fmt.Printf("%d name(s) in the array\n", len(names))
}

func TestLimitWithIter(t *testing.T) {
	var err error
	rd, err := os.Open("./petnames.txt")
	if err != nil {
		panic("Unable to open petnames file")
	}
	petNamesSeq := NewLineSeq(rd)
	newSeq := Limit(petNamesSeq, 5)
	expectedNames := []string{"AJ", "Abbey", "Abbie", "Abel", "Abigail"}
	var (
		i    int
		name string
	)
	for i, name = range IterWithIndex(newSeq) {
		expected := expectedNames[i]
		assert.Equal(t, expected, name)
	}
	assert.Equal(t, 4, i)
	assert.True(t, errors.Is(newSeq.Err(), io.EOF))
	assert.Nil(nil, petNamesSeq.lastErr)
}

func TestSkipSimple(t *testing.T) {
	strs := "foo\nbar\nbum\n"
	rd := strings.NewReader(strs)
	var sq Seq[string] = NewLineSeq(rd)
	sq = Skip(sq, 1)
	val, err := sq.Next()
	t.Logf("val=%v, err=%v\n", val, err)
	assert.Nil(t, err)
	assert.Equal(t, "bar", val)
}

func TestWhereNextThorough(t *testing.T) {
	var err error
	rd, err := os.Open("./petnames.txt")
	if err != nil {
		panic("Unable to open petnames file")
	}
	petNamesSeq := NewLineSeq(rd)
	filter := func(str string) bool { return strings.HasPrefix(str, "Ab") }
	newSeq := Where(petNamesSeq, filter)
	var name string
	// Abbey
	name, err = newSeq.Next()
	testWhere(t, "Abbey", name, nil, err)
	// Abbie
	name, err = newSeq.Next()
	testWhere(t, "Abbie", name, nil, err)
	// Abel
	name, err = newSeq.Next()
	testWhere(t, "Abel", name, nil, err)
	// Abigail
	name, err = newSeq.Next()
	testWhere(t, "Abigail", name, nil, err)
	// "", EOF
	name, err = newSeq.Next()
	testWhere(t, "", name, io.EOF, err)
}

func testWhere(t *testing.T, valExpected, val string, errExpected, err error) {
	t.Logf("name=%v err=%v\n", val, err)
	assert.Equal(t, errExpected, err)
	assert.Equal(t, valExpected, val)
}

func TestLimitAndWhereNextThorough(t *testing.T) {
	rd, err := os.Open("./petnames.txt")
	if err != nil {
		panic("Unable to open petnames file")
	}
	petNamesSeq := NewLineSeq(rd)
	filter := func(str string) bool { return strings.HasPrefix(str, "Ab") }
	newSeq := Where(petNamesSeq, filter)
	newerSeq := Limit(newSeq, 2)
	var name string
	// Abbey
	name, err = newerSeq.Next()
	testWhere(t, "Abbey", name, nil, err)
	// Abbie
	name, err = newerSeq.Next()
	testWhere(t, "Abbie", name, nil, err)
	// "", EOF
	name, err = newerSeq.Next()
	testWhere(t, "", name, io.EOF, err)
}

func TestWhereSkipLimit(t *testing.T) {
	rd, err := os.Open("./petnames.txt")
	if err != nil {
		panic("Unable to open petnames file")
	}
	sq := NewLineSeq(rd)
	var filter FilterFunc[string] = func(str string) bool { return strings.HasPrefix(str, "Ab") }
	sqWhere := Where(sq, filter)
	sqSkip := Skip(sqWhere, 3)
	sqLimit := Limit(sqSkip, 1)
	var name string
	// Abigail
	name, err = sqLimit.Next()
	t.Logf("name=%v err=%v\n", name, err)
	testWhere(t, "Abigail", name, nil, err)
	// "", EOF
	name, err = sqLimit.Next()
	t.Logf("name=%v err=%v\n", name, err)
	testWhere(t, "Abigail", name, io.EOF, err)
}
