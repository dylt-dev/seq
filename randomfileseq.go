package seq

import "os"

type RandomFileSeq struct {
	*HasErr[string]
	*IndexableFile
	path string
	sqLine *LineSeq
}

func (sq *RandomFileSeq) Count () int {
	n, err := sq.IndexableFile.Count()
	sq.SetErr(err)
	return n
}

func (sq *RandomFileSeq) Reset () (FiniteSeq[string], error) {
	var (f *os.File; err error)
	f, err = os.Open(sq.path)
	sq.SetErr(err)
	if err != nil { return sq, err }
	sq.sqLine = NewLineSeq(f)
	return sq, nil
}

func (sq *RandomFileSeq) Next () (string, error) {
	var err error
	if sq.sqLine == nil {
		var err error
		_, err = sq.Reset()
		if err != nil { return "", err}
	}
	var val string
	val, err = sq.sqLine.Next()
	sq.SetErr(err)
	return val, err
}