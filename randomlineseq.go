package seq

import "io"

// Seq for getting a random line from a file, without duplicates.
// 
// An internal FiniteLineCollection is used as the source of data. In general 
// random selection involves a fixed number of choices -- drawing one or more
// cards, rolling dice, picking names out of a hat, playing bingo, etc. 
// It's possible for a user to want a choice from an unbounded data source,
// but this necessarily means specifying a bound on the number of elements you
// want to consider, and choosing from among those elements.
type RandomLineSeq struct {
	*HasErr
	flc  FiniteLineCollection
	used map[string]struct{}
	shoe int
}

func NewRandomLineSeq(flc FiniteLineCollection, shoe int) *RandomLineSeq {
	return &RandomLineSeq{
		NewHasErr(),
		flc,
		map[string]struct{}{},
		shoe,
	}
}

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

