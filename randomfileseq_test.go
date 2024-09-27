package seq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomFile0 (t *testing.T) {
	path := "./petnames.txt"
	var sq *RandomFileSeq = &RandomFileSeq{
		path: path,
		sqLine: nil,
	}
	sq.HasErr = NewHasErr[string](sq)
	sq.HasCount = NewHasCount[string](sq)
	var _ SeqIndexable[string] = sq
	var (line string; err error)
	line, err = sq.Get(4)
	testNextOk(t, "Abigail", line, err)
}


func TestRandomFileCount (t *testing.T) {
	path := "./petnames.txt"
	var sq *RandomFileSeq = &RandomFileSeq{
		path: path,
		sqLine: nil,
	}
	sq.HasCount = NewHasCount(sq)
	sq.HasErr = NewHasErr[string](sq)
	var n int
	n = sq.Count()
	assert.Equal(t, 1000, n)
}
