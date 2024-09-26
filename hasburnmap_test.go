package seq

import "testing"

func TestHasBurnMap0 (t *testing.T) {
	var sqData Seq[int] = newArraySeq([]int{2, 2, 4, 4, 2, 3, 2, 3})
	var hbm *HasBurnMap[int] = NewHasBurnMap[int](3)
	var (val int; err error)
	val, err = hbm.AddFromSeq(sqData)
	t.Logf("val=%v err=%s isFull=%t %v\n", val, err, hbm.IsFull(), hbm.burned)
	val, err = hbm.AddFromSeq(sqData)
	t.Logf("val=%v err=%s isFull=%t %v\n", val, err, hbm.IsFull(), hbm.burned)
	val, err = hbm.AddFromSeq(sqData)
	t.Logf("val=%v err=%s isFull=%t %v\n", val, err, hbm.IsFull(), hbm.burned)
	val, err = hbm.AddFromSeq(sqData)
	t.Logf("val=%v err=%s isFull=%t %v\n", val, err, hbm.IsFull(), hbm.burned)
	val, err = hbm.AddFromSeq(sqData)
	t.Logf("val=%v err=%s isFull=%t %v\n", val, err, hbm.IsFull(), hbm.burned)
	t.Logf("IsFull: %t\n", hbm.IsFull())
}