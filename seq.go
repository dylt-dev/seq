package seq

import (
	"bufio"
	"errors"
	"io"
	"math/rand/v2"
	"strings"
)

type LoopFunc0 func() bool
type LoopFunc1[T comparable] func(arg T) bool
type LoopFunc2[T comparable] func(i int, arg T) bool
type IterFunc0 func(loopFunc LoopFunc0)
type IterFunc1[T comparable] func(loopFunc LoopFunc1[T])
type IterFunc2[T comparable] func(loopFunc LoopFunc2[T])
type NextFunc0 func() error
type NextFunc1[T comparable] func() (T, error)
type NextFunc2[T comparable] func() (int, T, error)

type Seq[T comparable] interface {
	Next() (T, error)
}

type HasErr struct {
	lastErr error
}

func NewHasErr() *HasErr {
	return &HasErr{lastErr: nil}
}

func (o *HasErr) Err() error {
	return o.lastErr
}

type HasIter[T any] struct {
	sq Seq[T]
}

func NewHasIter[T any] (sq Seq[T]) *HasIter[T] {
	return &HasIter[T]{sq}
}

func (o *HasIter[T]) Iter () IterFunc1[T] {
	return Iter1(o.sq)
}

func (o *HasIter[T]) IterWithIndex () IterFunc2[T] {
	return Iter2(o.sq)
}

func (o *HasIter[T]) IterNoArg () IterFunc0 {
	return Iter0(o.sq)
}



type HasPosition struct {
	lastPos int
	pos     int
}

func NewHasPosition() *HasPosition {
	return &HasPosition{lastPos: 0, pos: 0}
}

func (o *HasPosition) LastPosition() int {
	return o.lastPos
}

func (o *HasPosition) Position() int {
	return o.pos
}

func (o *HasPosition) Update(n int) int {
	o.lastPos = o.pos
	o.pos += n
	return o.pos
}

func Count[T comparable](seq Seq[T]) (int, error) {
	var err error
	var n int
	for n = 0; ; n++ {
		_, err := seq.Next()
		if err != nil {
			break
		}
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return n, err
}

type RuneSeq struct {
	*HasErr
	*HasPosition
	rd       io.Reader
	brd      bufio.Reader
	lastSize int
}

func NewRuneSeq(rd io.Reader) *RuneSeq {
	return &RuneSeq{
		HasErr:      NewHasErr(),
		HasPosition: NewHasPosition(),
		rd:          rd,
		brd:         *bufio.NewReaderSize(rd, 16),
		lastSize:    0,
	}
}

func (seq *RuneSeq) Next() (rune, error) {
	ru, size, err := seq.brd.ReadRune()
	seq.lastSize = size
	seq.lastErr = err
	seq.HasPosition.Update(size)
	return ru, err
}

type LineSeq struct {
	*HasErr
	*HasIter[string]
	*HasPosition
	rd      io.Reader
	runeSeq *RuneSeq
}

func NewLineSeq(rd io.Reader) *LineSeq {
	var sq *LineSeq = &LineSeq{
		HasErr:      NewHasErr(),
		HasPosition: NewHasPosition(),
		rd:          rd,
		runeSeq:     NewRuneSeq(rd),
	}
	// HasIter needs the Seq object so it needs special treatment
	sq.HasIter = NewHasIter(sq)
	return sq
}

func (seq *LineSeq) Err() error {
	return seq.runeSeq.Err()
}

func (seq *LineSeq) Next() (string, error) {
	b := strings.Builder{}
	for ru := range Iter1(seq.runeSeq) {
		if seq.runeSeq.lastSize > 0 {
			b.WriteRune(ru)
		}
		if ru == '\n' || errors.Is(seq.runeSeq.Err(), io.EOF) {
			break
		}
		if seq.runeSeq.Err() != nil {
			return "", seq.runeSeq.Err()
		}
	}
	seq.HasPosition.Update(b.Len())
	return strings.TrimRight(b.String(), "\n"), seq.runeSeq.Err()
}

type RandomLineSeq struct {
	*HasErr
	flc  FiniteLineCollection
	used map[string]struct{}
	shoe int
}

func NewRng () *rand.Rand {
	// t0 := uint64(time.Now().UnixNano())
	// t1 := uint64(time.Now().UnixNano())
	// src := rand.NewPCG(t0, t1)
	src := rand.NewPCG(rand.Uint64(), rand.Uint64())
	return rand.New(src)
}

func NewRandomLineSeq(flc FiniteLineCollection, shoe int) *RandomLineSeq {
	return &RandomLineSeq{
		NewHasErr(),
		flc,
		map[string]struct{}{},
		shoe,
	}
}

// func (seq *RandomLineSeq) Err() error {
// 	return seq.lastErr
// }

func (seq *RandomLineSeq) Next() (string, error) {
	flc := seq.flc
	n, err := flc.Count()
	if err != nil {
		return "", err
	}
	for {
		if seq.shoe+len(seq.used) >= n {
			return "", io.EOF
		}
		line, err := GetRandomLine(flc)
		if err != nil {
			seq.lastErr = err
			return "", err
		}
		_, hasKey := seq.used[line]
		if !hasKey {
			seq.used[line] = struct{}{}
			return line, nil
		}
	}
}

// func (seq *RandomLineSeq) Iter() RangeFunc1[string] {
// 	return CreateIterFromNext1(seq.Next)
// }

type SimpleRacer struct {
	HasErr
	speed float32
}

func NewSimpleRacer(speed float32) *SimpleRacer {
	return &SimpleRacer{speed: speed, HasErr: HasErr{lastErr: nil}}
}

func (racer *SimpleRacer) Next() (float32, error) {
	racer.lastErr = nil
	return racer.speed, racer.lastErr
}

// func CreateIterFromNext0(nextFunc NextFunc0) IterFunc0 {
// 	return func(loopFunc LoopFunc0) {
// 		for {
// 			err := nextFunc()
// 			if !loopFunc() || err != nil {
// 				break
// 			}
// 		}
// 	}
// }

// func CreateIterFromNext1[T any](nextFunc NextFunc1[T]) IterFunc1[T] {
// 	return func(loopFunc LoopFunc1[T]) {
// 		for {
// 			t, err := nextFunc()
// 			if !loopFunc(t) || err != nil {
// 				break
// 			}
// 		}
// 	}
// }

// func CreateIterFromNext2[T any](nextFunc NextFunc2[T]) IterFunc2[T] {
// 	fn := func(loopFunc LoopFunc2[T]) {
// 		for {
// 			int, string, err := nextFunc()
// 			if !loopFunc(int, string) || err != nil {
// 				break
// 			}
// 		}
// 	}
// 	return fn
// }

func Iter0[T comparable](seq Seq[T]) IterFunc0 {
	return func(loopFunc LoopFunc0) {
		for {
			_, err := seq.Next()
			if !loopFunc() || err != nil {
				break
			}
		}
	}
}

func Iter1[T comparable](seq Seq[T]) IterFunc1[T] {
	return func(loopFunc LoopFunc1[T]) {
		/*
		Next() == (non-empty, nil) => loopFunc(val)
		Next() == (non-empty, non-nil) => loopFunc(val)
		Next() == (empty, nil) => loopFunc(val)
		Next() -- (empty, non-nil) => break
		loopFunc == false => break
		*/
		var suchEmpty T = *new(T)
		if (t1 == *new(T)) { }
		for {
			t, err := seq.Next()
			if (t == *new(T) && err != nil) || !loopFunc(t) {
				break
			}
		}
	}
}

func Iter2[T comparable](seq Seq[T]) IterFunc2[T] {
	return func(loopFunc LoopFunc2[T]) {
		i := 0
		for {
			t, err := seq.Next()
			if !loopFunc(i, t) || err != nil {
				break
			}
			i++
		}
	}
}

type seqInner[T any] struct {
	i int
	limit int
}

func (sq *seqInner[T]) Next () (T, error) {
	if sq.i < sq.limit {
		sq.i++
		return sq.Next()
	}
	var t T
	return t, io.EOF
}

func Limit[T any] (seq Seq[T], limit int) Seq[T] {
	return &seqInner[T]{i: 0, limit: limit}
}