package seq

type FiniteSeq[T comparable] interface {
	Seq[T]
	Count () int
}