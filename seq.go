/* Package `seq` implements useful features for dealing with sequences of data. seq was Originally created
 * to make simplify creating Go Iterators, but it's grown a bit to provide stream-like features that are
 * familiar to other languages.
 */
package seq

import (
	"errors"
	"io"
)

// 'loop function' for 0-arg for loops (`for range iter`)
type LoopFunc0 func() bool
// 'loop function' for 1-arg for loops (`for el := range iter')
type LoopFunc1[T comparable] func(arg T) bool
// 'loop function' for 2-arg for loops (`for idx, el := range iter`)
type LoopFunc2[T comparable] func(i int, arg T) bool
// iterator for 0-arg for loops
type IterFunc0 func(loopFunc LoopFunc0)
// iterator for 1-arg for loops
type IterFunc1[T comparable] func(loopFunc LoopFunc1[T])
// iterator for 2 arg for loops
type IterFunc2[T comparable] func(loopFunc LoopFunc2[T])
type NextFunc0 func() error
type NextFunc1[T comparable] func() (T, error)
type NextFunc2[T comparable] func() (int, T, error)

// Seq is the most fundamental type of the `seq` package. In true Go fashion it only has a single method,
// `Next()`. `Next()` returns the next element in the sequences, along with any error that may have occured.
// This provides an opportunity for handling errors that is absent from Go Iterators.
//
// In a typical use case, `Seq` will be implemented by a struct of the desired base type, using a specific data
// source. `Next()` will get a single element from this data source, keeping track of the current position if
// needed. `Next()` indicates completion by returning *new(T), io.EOF, where `*new(T)` represents an empty value for
// type T, eg 0, nil, of "". Go currently has no better idiom for representing a generic empty value. Consumers of
// a Seq, eg iterators, can check for return values of (*new*T(), nil) to know when to stop calling Next() and finish
// up. 
//
// structs that implement `Seq` might want to make other information available about the last operation or about
// the aggregate use of the `Seq`. File-based `Seq`s might track file position ... Network-based `Seq`s might track
// total bytes received ... it's entirely up to the creator. Users who create Iterators from `Seq`s can check `Seq`
// custom properties as needed. The only requirement is that the struct's Seq.Next()` method updates those custom 
// properties as appropriate.
type Seq[T comparable] interface {
	// Get the next element in the sequence, and an error or nil
	Next() (T, error)
}


// Next() through the entire sequence until exhaustion to count the number of elements
//
// Warning 1 - Depending on the underlying datasource, this might consume and discard data that you did not want consumed and discarded. For data structures and local files, this is probably ok. For remote data and service resposnes, this could be a problem.
// Warninf 2 - Not all sequences end. Some might be an infinite sequence of data. If so, Count() will not terminate.
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

func IterNoArg[T comparable](seq Seq[T]) IterFunc0 {
	return func(loopFunc LoopFunc0) {
		for {
			t, err := seq.Next()
			if (t == *new(T) && err != nil) || !loopFunc() {
				break
			}
		}
	}
}

// The function that this was all built to support. `Iter()` takes a sequence and returns
// an iterator ... specifically it returns the func (func (val T) bool) flavor of iterator,
// designed to work with for .. range loops that use the value but not the index of each
// element.
func Iter[T comparable](seq Seq[T]) IterFunc1[T] {
	return func(loopFunc LoopFunc1[T]) {
		/*
			Next() == (non-empty, nil) => loopFunc(val)
			Next() == (non-empty, non-nil) => loopFunc(val)
			Next() == (empty, nil) => loopFunc(val)
			Next() -- (empty, non-nil) => break
			loopFunc == false => break
			(*new)T is 'the zero value irrespective of generic type' in Go
		*/
		for {
			t, err := seq.Next()
			if (t == *new(T) && err != nil) || !loopFunc(t) || err != nil {
				break
			}
		}
	}
}

func IterWithIndex[T comparable](seq Seq[T]) IterFunc2[T] {
	return func(loopFunc LoopFunc2[T]) {
		i := 0
		for {
			t, err := seq.Next()
			if (t == *new(T) && err != nil) || !loopFunc(i, t) {
				break
			}
			i++
		}
	}
}

type seqLimit[T comparable] struct {
	*HasErr[T]
	sqInner Seq[T]
	i       int
	limit   int
}

func NewSeqLimit[T comparable] (sqInner Seq[T], limit int) *seqLimit[T]{
	var sq *seqLimit[T] = &seqLimit[T] {
		sqInner: sqInner,
		i: 0,
		limit: limit,
	}
	sq.HasErr = NewHasErr(sq)
	return sq
}

func (sq *seqLimit[T]) Next() (T, error) {
	if sq.i < sq.limit {
		sq.i++
		tInner, errInner := sq.sqInner.Next()
		// If this is the last iteration, explcitly set err to EOF
		if errInner == nil && sq.i == sq.limit {
			errInner = io.EOF
		}
		sq.lastErr = errInner
		return tInner, errInner
	}
	sq.lastErr = io.EOF
	return *new(T), io.EOF
}

func Limit[T comparable](sqInner Seq[T], limit int) *seqLimit[T] {
	return NewSeqLimit(sqInner, limit)
}

type FilterFunc[T comparable] func(T) bool

type seqWhere[T comparable] struct {
	*HasErr[T]
	sqInner Seq[T]
	filter  FilterFunc[T]
}

func NewSeqWhereWrapper[T comparable](sqInner Seq[T], filter FilterFunc[T]) *seqWhere[T] {
	var sq *seqWhere[T] = &seqWhere[T]{
		sqInner: sqInner,
		filter: filter,
	}
	sq.HasErr = NewHasErr(sq)
	return sq
}

func (sq *seqWhere[T]) Next() (T, error) {
	for {
		next, err := sq.sqInner.Next()
		if next == *new(T) && err != nil {
			return next, err
		}
		if sq.filter(next) {
			return next, err
		}
	}
}

func Where[T comparable](sqInner Seq[T], filter FilterFunc[T]) *seqWhere[T] {
	return NewSeqWhereWrapper(sqInner, filter)
}

type seqSkip[T comparable] struct {
	*HasErr[T]
	sqInner   Seq[T]
	toSkip    int
	isSkipped bool
}

func NewSeqSkipWrapper[T comparable](sqInner Seq[T], toSkip int) *seqSkip[T] {
	var sq *seqSkip[T] = &seqSkip[T]{
		sqInner: sqInner,
		toSkip: toSkip,
		isSkipped: false,
	}
	sq.HasErr = NewHasErr(sq)
	return sq
}

func (sq *seqSkip[T]) Next() (T, error) {
	//  On first call to Next() we do our skip step
	if !sq.isSkipped {
		sq.isSkipped = true
		// skip step: Skip over the specified # of elements
		for range sq.toSkip {
			_, err := sq.sqInner.Next()
			// If there's an error during the skip step, terminate early
			if err != nil {
				sq.lastErr = err
				return *new(T), err
			}
		}
	}
	// We've done the skip step, so now we just delegate to sqInner.Next()
	val, err := sq.sqInner.Next()
	sq.lastErr = err
	return val, err
}

func Skip[T comparable](sq Seq[T], toSkip int) *seqSkip[T] {
	return NewSeqSkipWrapper(sq, toSkip)
}


type SeqWithErr[T comparable] interface {
	Seq[T]
	Err () error
	SetErr (err error) SeqWithErr[T]
}

type SeqIndexable[T comparable] interface {
	FiniteSeq[T]
	Get (i int) (T, error)
}