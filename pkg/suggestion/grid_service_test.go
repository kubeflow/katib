package suggestion

import (
	"testing"
)

func TestAllocInt(t *testing.T) {
	s := &GridSuggestService{}

	min1 := 1
	max1 := 2
	reqnum1 := 1

	exp1 := "1"
	rtn1 := s.allocInt(min1, max1, reqnum1)

	if rtn1[0] != exp1 {
		t.Errorf("expected [%v], but %v is returned",exp1, rtn1)
	}

	min2 := 1
	max2 := 9
	reqnum2 := 5

	exp2 := []string{"1", "3", "5", "7", "9"}
	rtn2 := s.allocInt(min2, max2, reqnum2)

	for i := 0; i < 5; i++ {
		if rtn2[i] != exp2[i] {
			t.Errorf("expected Array[%v] = %v, but %v is returned", i, exp2[i], rtn2[i])
		}
	}
}

func TestAllocFloat(t *testing.T) {
	s := &GridSuggestService{}

	min1 := 1.0
	max1 := 2.0
	reqnum1 := 1

	exp1 := "1.0000"
	rtn1 := s.allocFloat(min1, max1, reqnum1)

	if rtn1[0] != exp1 {
		t.Errorf("expected [%v], but %v is returned",exp1, rtn1)
	}

	min2 := 1.0
	max2 := 9.0
	reqnum2 := 5

	exp2 := []string{"1.0000", "3.0000", "5.0000", "7.0000", "9.0000"}
	rtn2 := s.allocFloat(min2, max2, reqnum2)

	for i := 0; i < 5; i++ {
		if rtn2[i] != exp2[i] {
			t.Errorf("expected Array[%v] = %v, but %v is returned", i, exp2[i], rtn2[i])
		}
	}
}

func TestAllocCat(t *testing.T) {
	s := &GridSuggestService{}

	list := []string{"alpha", "beta", "gamma"}
	reqnum1 := 1

	exp1 := []string{"alpha"}
	rtn1 := s.allocCat(list, reqnum1)

	if rtn1[0] != exp1[0] {
		t.Errorf("exptected %v, but %v", rtn1, exp1)
	}

	reqnum2 := 5
	exp2 := []string{"alpha", "alpha", "beta", "beta", "gamma"}
	rtn2 := s.allocCat(list, reqnum2)

	for i := 0; i < 5; i++ {
		if rtn2[i] != exp2[i] {
			t.Errorf("expected Array[%v] = %v, but %v is returned", i, exp2[i], rtn2[i])
		}
	}
}
