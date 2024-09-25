package seq

import (
	"errors"
	"io"
	"strings"
)

// Seq for consuming a Reader line by line
//
// Add-ons: HasErr, HasIter, HasPosition
//
// LineSeq uses a RuneSeq internally to consume a Reader rune-by-rune in order to form lines.
// The main reason for this was simply to indirectly test RuneSeq. A nice side effect is that
// using RuneSeq and building lines ourselves is potentially less memory-intensive than using
// bufio.ReadString('\n'), since the latter creates internal buffers of a size we can't control.
// This is unlikely to be that big a deal in reality, but it's nice to know.
type LineSeq struct {
	*HasErr
	*HasIter[string]
	*HasPosition
	rd      io.Reader
	runeSeq *RuneSeq
}

// C'tor last function
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

// LineSeq delegates its actual reading to RuneSeq, so Err() delegates and
// returns the last RuneSeq error. (@todo maybe we can/should drop *HasErr?)
func (seq *LineSeq) Err() error {
	return seq.runeSeq.Err()
}

// Read runes until '\n' or EOF is reached.
//
// After EOF is reached, all further calls to Next() will return ("", io.EOF)
func (seq *LineSeq) Next() (string, error) {
	b := strings.Builder{}
	for ru := range Iter(seq.runeSeq) {
		// Make sure that the RuneSeq actually read something
		if seq.runeSeq.lastSize > 0 {
			b.WriteRune(ru)
		}
		// Terminate either on '\n' or EOF. This correctly files that don't end in '\n'.
		if ru == '\n' || errors.Is(seq.runeSeq.Err(), io.EOF) {
			break
		}
		if seq.runeSeq.Err() != nil {
			return "", seq.runeSeq.Err()
		}
	}
	// A string was succesfully read, so update the position variables
	seq.HasPosition.Update(b.Len())
	// Remove the '\n'. Technically this removes all trailing '\n's but since we stop at '\n' that's not a concern.
	str := strings.TrimRight(b.String(), "\n")
	return str, seq.runeSeq.Err()
}
