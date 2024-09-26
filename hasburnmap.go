package seq

import (
	"errors"
)

var ErrMapFull error = errors.New("ErrMapFull")

type HasBurnMap[T comparable] struct {
	burned map[T]struct{}
	maxSize int
}

func NewHasBurnMap[T comparable] (maxSize int) *HasBurnMap[T] {
	return &HasBurnMap[T]{
		burned: make(map[T] struct{}),
		maxSize: maxSize,
	}
}

func (o *HasBurnMap[T]) AddFromSeq (sq Seq[T]) (T, error) {
	var (val T; err error)
	for !o.IsFull() {
		val, err = sq.Next()
		if val == *new(T) && err != nil {
			break
		}
		_, hasKey := o.burned[val]
		if !hasKey {
			o.burned[val] = struct{}{}
			break
		}
	}
	if o.IsFull() {
		err = ErrMapFull
	}
	return val, err
}

func (o *HasBurnMap[T]) Capacity () int {
	return max(0, o.maxSize - len(o.burned))
}

func (o *HasBurnMap[T]) IsFull () bool {
	return (o.Capacity() <= 0)
}
