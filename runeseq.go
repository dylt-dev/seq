package seq

import (
	"bufio"
	"io"
)

// Seq for expressing a file as a sequence of runes. The Seq that started it all!
// A bufio.Reader() serves as the underlying data source. bufio.Reader.ReadRune() does
// most of the heavy lifting.
//
// Runes can be more than one character, so when ReadRune() returns a rune, it also returns
// the number of bytes in the rune. Using an iterator, this information would be lost. So RuneSeq
// makes it available via the `LastSize()` method.
type RuneSeq struct {
	*HasErr
	*HasIter[rune]
	*HasPosition
	rd       io.Reader
	brd      bufio.Reader
	lastSize int
}

// C'tor function
func NewRuneSeq(rd io.Reader) *RuneSeq {
	sq := &RuneSeq{
		HasErr:      NewHasErr(),
		HasPosition: NewHasPosition(),
		rd:          rd,
		brd:         *bufio.NewReaderSize(rd, 16),
		lastSize:    0,
	}
	// HasIter needs the Seq object so it needs special treatment
	sq.HasIter = NewHasIter(sq)
	return sq
}

// Return the next rune in the sequence.
// Extra fields: last error, last/current position, last rune size in bytes
// ReadRune() returns (0, io.EOF) upon exhaustion, so Next() doesn't have to do
// anything special to detect when there's more data. ReadRune() handles it.
func (seq *RuneSeq) Next() (rune, error) {
	ru, size, err := seq.brd.ReadRune()
	seq.lastSize = size
	seq.lastErr = err
	seq.HasPosition.Update(size)
	return ru, err
}
