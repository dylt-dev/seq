/* Package `seq` implements useful features for dealing with sequences of data. seq was Originally created
 * to make simplify creating Go Iterators, but it's grown a bit to provide stream-like features that are
 * familiar to other languages.
 */
package seq

import (
	"bufio"
	"errors"
	"io"
	"math/rand/v2"
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

// Seq Add-on for tracking the last error received by `Next()`. Can be used to check if the Seq completed normally (io.EOF),
// or if some other error happened.
// Proper usage requires including HasErr or *HasErr as an embedded field, and then making sure Next() updates lastErr with
// any errors that occur or nil if Next() executes successfully.
type HasErr struct {
	lastErr error
}

// C'tor function
func NewHasErr() *HasErr {
	return &HasErr{lastErr: nil}
}

// Return the last error
func (o *HasErr) Err() error {
	return o.lastErr
}

// Seq Add-on to use Iter(), IterWithIndex(), and IterNoArg() as methods, instead of global functions.
// HasIter's sq field typically represnts the Seq that embeds HasIter. This means Seq's embedding HasIter
// cannot initialize HasIter when the Seq is initialized, beause the Seq doesn't exist yet. Instead, create
// the Seq first, then explicitly create a HasIter and initialize it with a pointer to the new Seq, etc.
//
//  var sq *MyNewSeq = &MyNewSeq{}
//	var hasIter *HasIter = NewHasIter{sq}
//  sq.HasIter = hasIter
//
// If this seems like too much work to support IterXXX() functions as methods, you can just skip HasIter, but
// some users really like using methods. Up to you.
type HasIter[T comparable] struct {
	sq Seq[T]
}

// C'tor function.
func NewHasIter[T comparable](sq Seq[T]) *HasIter[T] {
	return &HasIter[T]{sq}
}

// Return an iterator. Usage: `for val := range sq.Iter()`
func (o *HasIter[T]) Iter() IterFunc1[T] {
	return Iter(o.sq)
}

// Return an (index, val) iterator. Usage: `for idx, val := range sq.IterWithIndex()`
func (o *HasIter[T]) IterWithIndex() IterFunc2[T] {
	return IterWithIndex(o.sq)
}

// Return a no-arg iterator. Usage: `for range := sq.IterNoArg()`
func (o *HasIter[T]) IterNoArg() IterFunc0 {
	return IterNoArg(o.sq)
}

// Add-on for tracking position of the underlying data source. Example: io.Reader character position
// for a Seq of lines or tokens. Two positions are available: the last element returned by Next(), and
// the current position that Next() will use the next time it's called.
type HasPosition struct {
	lastPos int
	pos     int
}

// C'tor function
func NewHasPosition() *HasPosition {
	return &HasPosition{lastPos: 0, pos: 0}
}

// Position of last returned element
func (o *HasPosition) LastPosition() int {
	return o.lastPos
}

// Position of next element
func (o *HasPosition) Position() int {
	return o.pos
}

// Rotate Position to LastPosition, increment new position by n, and return the previous position
func (o *HasPosition) Update(n int) int {
	o.lastPos = o.pos
	o.pos += n
	return o.pos
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
	*HasErr
	sqInner Seq[T]
	i       int
	limit   int
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
	return &seqLimit[T]{NewHasErr(), sqInner, 0, limit}
}

type FilterFunc[T comparable] func(T) bool

type seqWhere[T comparable] struct {
	*HasErr
	sqInner Seq[T]
	filter  FilterFunc[T]
}

func NewSeqWhereWrapper[T comparable](sqInner Seq[T], filter FilterFunc[T]) *seqWhere[T] {
	return &seqWhere[T]{NewHasErr(), sqInner, filter}
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
	*HasErr
	sqInner   Seq[T]
	toSkip    int
	isSkipped bool
}

func NewSeqSkipWrapper[T comparable](sqInner Seq[T], toSkip int) *seqSkip[T] {
	return &seqSkip[T]{NewHasErr(), sqInner, toSkip, false}
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
