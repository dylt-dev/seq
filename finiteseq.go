package seq

type FiniteSeq[T comparable] interface {
	SeqWithErr[T]
	Count () int
	Reset () (FiniteSeq[T], error)
}