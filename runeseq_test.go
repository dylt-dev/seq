package seq

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)
func TestRuneSeq0(t *testing.T) {
	str := TS0
	sr := strings.NewReader(str)
	sq := NewRuneSeq(sr)
	for i, ru := range IterWithIndex(sq) {
		t.Logf("i=%d ru=%c\n", i, ru)
	}
	testEof(t, sq)
}

func TestRuneSeqCount(t *testing.T) {
	str := TS0
	sr := strings.NewReader(str)
	sq := NewRuneSeq(sr)
	n, err := Count(sq)
	assert.Equal(t, len(str), n)
	assert.Nil(t, err)
	testEof(t, sq)
}

