package seq

import (
	_ "embed"
	"errors"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed iter-test-0.txt
var TS0 string

/*
func TestSimpleRacer(t *testing.T) {
	speed := float32(2.27)
	racer := NewSimpleRacer(speed)
	for i, dx := range IterWithIndex(racer) {
		t.Logf("%f", dx)
		time.Sleep(10 * time.Millisecond)
		if i >= 10 {
			break
		}
	}
}
*/


func TestLimitNextThorough(t *testing.T) {
	var err error
	rd, err := os.Open("./petnames.txt")
	if err != nil {
		panic("Unable to open petnames file")
	}
	petNamesSeq := NewLineSeq(rd)
	// names := []string{}
	var newSeq *seqLimit[string] = Limit(petNamesSeq, 5)
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
	// 4: "Abigail", EOF
	name, err = newSeq.Next()
	testNext(t, "Abigail", name, io.EOF, err)
}

func TestLimitWithIter(t *testing.T) {
	var err error
	rd, err := os.Open("./petnames.txt")
	if err != nil {
		panic("Unable to open petnames file")
	}
	var sq Seq[string] = NewLineSeq(rd)
	var sqLimit *seqLimit[string] = Limit(sq, 5)
	expectedNames := []string{"AJ", "Abbey", "Abbie", "Abel", "Abigail"}
	var (i int; name string)
	for i, name = range IterWithIndex(sqLimit) {
		expected := expectedNames[i]
		assert.Equal(t, expected, name)
	}
	assert.Equal(t, 4, i)
	assert.True(t, errors.Is(sqLimit.Err(), io.EOF))
	// EOF
	testEof(t, sqLimit)
}

func TestSkipSimple(t *testing.T) {
	strs := "foo\nbar\nbum\n"
	rd := strings.NewReader(strs)
	var sq Seq[string] = NewLineSeq(rd)
	sq = Skip(sq, 1)
	var (val string; err error)
	// bar
	val, err = sq.Next()
	testNextOk(t, "bar", val, err)
	// bum
	val, err = sq.Next()
	testNextOk(t, "bum", val, err)
	// "", EOF
	val, err = sq.Next()
	testNextEof(t, val, err)
	// Confirm EOF
	testEof(t, sq)
}

func TestWhereNextThorough(t *testing.T) {
	var err error
	rd, err := os.Open("./petnames.txt")
	if err != nil {
		panic("Unable to open petnames file")
	}
	sq := NewLineSeq(rd)
	filter := func(str string) bool { return strings.HasPrefix(str, "Ab") }
	sqWhere := Where(sq, filter)
	var name string
	// Abbey
	name, err = sqWhere.Next()
	testNext(t, "Abbey", name, nil, err)
	// Abbie
	name, err = sqWhere.Next()
	testNext(t, "Abbie", name, nil, err)
	// Abel
	name, err = sqWhere.Next()
	testNext(t, "Abel", name, nil, err)
	// Abigail
	name, err = sqWhere.Next()
	testNext(t, "Abigail", name, nil, err)
	// "", EOF
	name, err = sqWhere.Next()
	testNext(t, "", name, io.EOF, err)
	// Confirm EOF
	testEof(t, sqWhere)
}

func TestLimitAndWhereNextThorough(t *testing.T) {
	rd, err := os.Open("./petnames.txt")
	if err != nil {
		panic("Unable to open petnames file")
	}
	sq := NewLineSeq(rd)
	filter := func(str string) bool { return strings.HasPrefix(str, "Ab") }
	sqWhere := Where(sq, filter)
	sqLimit := Limit(sqWhere, 2)
	var name string
	// Abbey
	name, err = sqLimit.Next()
	testNext(t, "Abbey", name, nil, err)
	// Abbie, EOF
	name, err = sqLimit.Next()
	testNext(t, "Abbie", name, io.EOF, err)
	// Confirm EOF
	testEof(t, sqLimit)
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
	// Abigail, EOF
	name, err = sqLimit.Next()
	t.Logf("name=%v err=%v\n", name, err)
	testNext(t, "Abigail", name, io.EOF, err)
	// "", EOF
	name, err = sqLimit.Next()
	t.Logf("name=%v err=%v\n", name, err)
	testNextEof(t, name, err)
	// Confirm EOF
	testEof(t, sqLimit)
}


func testEof[U comparable] (t *testing.T, sq Seq[U]) {
	var (val U; err error)
	val, err = sq.Next()
	t.Logf("name=%v err=%v\n", val, err)
	assert.Equal(t, *new(U), val)
	assert.True(t, errors.Is(err, io.EOF))
}


func testNext[U comparable](t *testing.T, valExpected, val U, errExpected, err error) {
	t.Logf("name=%v err=%v\n", val, err)
	assert.Equal(t, errExpected, err)
	assert.Equal(t, valExpected, val)
}

func testNextEof[U comparable](t *testing.T, val U, err error) {
	testNext(t, *new(U), val, io.EOF, err)
}

func testNextOk[U comparable](t *testing.T, valExpected, val U, err error) {
	testNext(t, valExpected, val, nil, err)
}

// This type doesn't have much of a purpose outside of testing
type arraySeq[T comparable] struct {
	*HasCount[T]
	*HasErr[T]
	*HasIter[T]
	data []T
	i int
}

func newArraySeq[T comparable] (data []T) *arraySeq[T] {
	var sq *arraySeq[T] = &arraySeq[T]{
		data: data,
		i: 0,
	}
	sq.HasErr = NewHasErr(sq)
	sq.HasCount = NewHasCount(sq)
	sq.HasIter = NewHasIter(sq)
	return sq
}

var ErrIndexOutOfBounds error = errors.New("ErrIndexOutOfBounds")

func (sq *arraySeq[T]) Get (i int) (T, error) {
	if i >= len(sq.data) {
		return *new(T), ErrIndexOutOfBounds
	} 
	return sq.data[i], nil
}

func (sq *arraySeq[T]) Next () (val T, err error) {
	if sq.i < len(sq.data) {
		val = sq.data[sq.i]
		err = nil
		sq.i++
		return
	}
	val = *new(T)
	err = io.EOF
	return
}

func (sq *arraySeq[T]) Reset () (FiniteSeq[T], error)  {
	sq.i = 0
	return sq, nil
}

func TestArraySeqNextThorough (t *testing.T) {
	var sq *arraySeq[int] = newArraySeq[int]([]int{2, 4, 3, 5, 1})
	var (val int; err error)
	val, err = sq.Next()
	assert.Nil(t, err)
	assert.Equal(t, 2, val)
	val, err = sq.Next()
	assert.Nil(t, err)
	assert.Equal(t, 4, val)
	val, err = sq.Next()
	assert.Nil(t, err)
	assert.Equal(t, 3, val)
	val, err = sq.Next()
	assert.Nil(t, err)
	assert.Equal(t, 5, val)
	val, err = sq.Next()
	assert.Nil(t, err)
	assert.Equal(t, 1, val)
	val, err = sq.Next()
	assert.NotNil(t, err)
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, 0, val)
}

func TestArraySeqCount (t *testing.T) {
	var sq *arraySeq[int] = newArraySeq[int]([]int{2, 4, 3, 5, 1})
	assert.Equal(t, 5, sq.Count())
	assert.Equal(t, nil, sq.Err())
	testEof(t, sq)
	sq.Reset()
	val, err := sq.Next()
	testNextOk(t, 2, val, err)
}

func TestIndexable0 (t *testing.T) {
	data := []int{2, 4, 3, 5, 1}
	var sq *arraySeq[int] = newArraySeq[int](data)
	var (val int; err error)
	// data[4] = 1
	val, err = sq.Get(4)
	assert.Nil(t, err)
	assert.Equal(t, 1, val)
	// data[0] = 2
	val, err = sq.Get(0)
	assert.Nil(t, err)
	assert.Equal(t, 2, val)
	// data[2] = 3
	val, err = sq.Get(2)
	assert.Nil(t, err)
	assert.Equal(t, 3, val)
	// data[1] = 4
	val, err = sq.Get(1)
	assert.Nil(t, err)
	assert.Equal(t, 4, val)
	// data[3] = 5
	val, err = sq.Get(3)
	assert.Nil(t, err)
	assert.Equal(t, 5, val)
	// data[666] = 0, Err
	val, err = sq.Get(666)
	assert.NotNil(t, err)
	assert.True(t, errors.Is(err, ErrIndexOutOfBounds))
	assert.Equal(t, 0, val)
}