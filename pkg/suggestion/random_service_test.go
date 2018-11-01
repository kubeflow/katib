package suggestion
import (
	"testing"
)
func TestDoubleRandom(t *testing.T) {
	s := &RandomSuggestService{}
	min1 := 1.0
	max1 := 1.0
	rtn1 := s.DoubleRandom(min1, max1)
	exp1 := 1.0
	if rtn1 != exp1 {
		t.Errorf("actual %v expected %v", rtn1, exp1)
	}
	min2 := 1.0
	max2 := 10.0
	var rtnA,rtnB float64
	for i := 0; i <= 10; i++ {
		rtnA = s.DoubleRandom(min2, max2)
		rtnB = s.DoubleRandom(min2, max2)
		if rtnA != rtnB {
			break
		} else if i == 10 {
			t.Errorf("different value is expected, but same value is returned %v = %v", rtnA, rtnB)
		}
	}
}
func TestIntRandom(t *testing.T) {
	s := &RandomSuggestService{}
	min := 1
	max := 10
	var rtnA,rtnB int
	for i := 0; i <= 10; i++ {
		rtnA = s.IntRandom(min, max)
		rtnB = s.IntRandom(min, max)

		if rtnA != rtnB {
			break
		} else if i == 10 {
			t.Errorf("different value is expected, but same value is returned %v = %v", rtnA, rtnB)
		}
	}
}
