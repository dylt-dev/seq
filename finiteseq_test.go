package seq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFiniteSeq0 (t *testing.T) {
	var data []int = []int{2, 4, 3, 5, 1}
	var sq FiniteSeq[int] = newArraySeq(data)
	assert.Equal(t, len(data), sq.Count())
	sq.Reset()
	val, err := sq.Next()
	assert.Nil(t, err)
	assert.Equal(t, 2, val)
}

func TestFiniteSeqIter (t *testing.T) {
	var data []int = []int{2, 4, 3, 5, 1}
	var sq *arraySeq[int] = newArraySeq(data)
	expectedValues := []int{2, 4, 3, 5, 1}
	expectedErrs := []error{nil, nil, nil, nil, nil}
	for i, val := range sq.IterWithIndex() {
		assert.Equal(t, expectedErrs[i], sq.Err())
		assert.Equal(t, expectedValues[i], val)
	}
	testEof(t, sq)
}