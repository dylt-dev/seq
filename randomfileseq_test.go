package seq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomFile0(t *testing.T) {
	path := "./petnames.txt"
	var sq *RandomFileSeq = &RandomFileSeq{
		path:   path,
		sqLine: nil,
	}
	sq.HasErr = NewHasErr[string](sq)
	sq.IndexableFile = &IndexableFile{path: path}
	var _ Indexable[string] = sq.IndexableFile
	var (
		line string
		err  error
	)
	line, err = sq.Get(4)
	testNextOk(t, "Abigail", line, err)
}

func TestRandomFileCount(t *testing.T) {
	path := "./petnames.txt"
	var sq *RandomFileSeq = &RandomFileSeq{
		path:   path,
		sqLine: nil,
	}
	sq.HasErr = NewHasErr[string](sq)
	sq.IndexableFile = &IndexableFile{path: path}
	var n int
	n = sq.Count()
	assert.Equal(t, 1000, n)
}
