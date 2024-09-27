package seq

import "os"

type RandomFileSeq struct {
	*HasCount[string]
	*HasErr[string]
	path string
	sqLine *LineSeq
}

func (sq *RandomFileSeq) Reset () (FiniteSeq[string], error) {
	var (f *os.File; err error)
	f, err = os.Open(sq.path)
	sq.SetErr(err)
	if err != nil { return nil, err }
	sq.sqLine = NewLineSeq(f)
	return sq, nil
}

func (sq *RandomFileSeq) Get (i int) (string, error) {
	var err error
	if sq.sqLine == nil {
		_, err = sq.Reset()
		if err != nil { return "", err}
	}
	var val string
	for range i {
		val, err = sq.sqLine.Next()
		sq.SetErr(err)
		if val == "" && err != nil {
			break
		}
	}
	val, err = sq.sqLine.Next()
	return val, err
}

func (sq *RandomFileSeq) Next () (string, error) {
	var err error
	if sq.sqLine == nil {
		var err error
		_, err = sq.Reset()
		sq.SetErr(err)
		if err != nil { return "", err}
	}
	var val string
	val, err = sq.sqLine.Next()
	sq.SetErr(err)
	return val, err
}