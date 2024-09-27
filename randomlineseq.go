package seq

import "io"

type RandomLineSeq struct {
	*HasErr[string]
	flc  FiniteLineCollection
	used map[string]struct{}
	shoe int
}

func NewRandomLineSeq(flc FiniteLineCollection, shoe int) *RandomLineSeq {
	var sq *RandomLineSeq = &RandomLineSeq{
		flc: flc,
		used: map[string]struct{}{},
		shoe: shoe,
	}
	sq.HasErr = NewHasErr(sq)
	return sq
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

