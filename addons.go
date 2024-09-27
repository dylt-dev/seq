package seq

import (
	"errors"
	"io"
)

// Seq Add-on for returning the Count() of the sequence.
// This implementation is naive. It simple Nexts() all the way to the end
// of the sequence. Not only is that potentially slow, but it consumes the
// entire sequence. For some sequences that's fine. For others, it will be
// better to skip this Add-on and code a Seq-specific Count().
//
// Implements: FiniteSeq
type HasCount[T comparable] struct {
	sq FiniteSeq[T]
}

func NewHasCount[T comparable] (sq FiniteSeq[T]) *HasCount[T] {
	return &HasCount[T]{sq}
}

func (o *HasCount[T]) Count () int {
	var (i int=0; val T; err error)
	for {
		val, err = o.sq.Next()
		if val != *new(T) {
			i++
		}
		if err != nil {
			break
		} 
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	o.sq.SetErr(err)
	return i
}

// Seq Add-on for tracking the last error received by `Next()`. Can be used to check if the Seq completed normally (io.EOF),
// or if some other error happened.
// Proper usage requires including HasErr or *HasErr as an embedded field, and then making sure Next() updates lastErr with
// any errors that occur or nil if Next() executes successfully.
type HasErr[T comparable] struct {
	sq SeqWithErr[T]
	lastErr error
}

// C'tor function
func NewHasErr[T comparable](sq SeqWithErr[T]) *HasErr[T] {
	return &HasErr[T]{sq: sq, lastErr: nil}
}

// Return the last error
func (o *HasErr[T]) Err() error {
	return o.lastErr
}

// Set the last error
func (o *HasErr[T]) SetErr (err error) SeqWithErr[T] {
	o.lastErr = err
	return o.sq
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

