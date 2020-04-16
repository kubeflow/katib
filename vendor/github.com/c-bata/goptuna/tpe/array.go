package tpe

import (
	"sort"
)

func ones1d(size int) []float64 {
	ones := make([]float64, size)
	for i := 0; i < size; i++ {
		ones[i] = 1
	}
	return ones
}

func linspace(start, stop float64, num int, endPoint bool) []float64 {
	step := 0.
	if endPoint {
		if num == 1 {
			return []float64{start}
		}
		step = (stop - start) / float64(num-1)
	} else {
		if num == 0 {
			return []float64{}
		}
		step = (stop - start) / float64(num)
	}
	r := make([]float64, num, num)
	for i := 0; i < num; i++ {
		r[i] = start + float64(i)*step
	}
	return r
}

func choice(array []float64, idxs []int) []float64 {
	results := make([]float64, len(idxs))
	for i, idx := range idxs {
		results[i] = array[idx]
	}
	return results
}

func location(array []float64, key float64) int {
	i := 0
	size := len(array)
	for {
		mid := (i + size) / 2
		if i == size {
			break
		}
		if array[mid] < key {
			i = mid + 1
		} else {
			size = mid
		}
	}
	return i
}

func searchsorted(array, values []float64) []int {
	var indexes []int
	for _, val := range values {
		indexes = append(indexes, location(array, val))
	}
	return indexes
}

func bincount(x []int, weights []float64, minlength int) []float64 {
	// Count the number of occurrences of each value in array of non-negative ints.
	// https://docs.scipy.org/doc/numpy/reference/generated/numpy.bincount.html
	counts := make([]float64, minlength)
	for i := range x {
		if x[i] > len(counts)-1 {
			for j := len(counts) - 1; j < x[i]; j++ {
				counts = append(counts, 0)
			}
		}
		if x[i] > len(weights)-1 {
			counts[x[i]]++
		} else {
			counts[x[i]] += weights[x[i]]
		}
	}
	return counts
}

func clip(array []float64, min, max float64) {
	for i := range array {
		if array[i] < min {
			array[i] = min
		} else if array[i] > max {
			array[i] = max
		}
	}
}

func argSort2d(lossVals [][2]float64) []int {
	type sortable struct {
		index   int
		lossVal [2]float64
	}
	x := make([]sortable, len(lossVals))
	for i := 0; i < len(lossVals); i++ {
		x[i] = sortable{
			index:   i,
			lossVal: lossVals[i],
		}
	}

	sort.SliceStable(x, func(i, j int) bool {
		if x[i].lossVal[0] == x[j].lossVal[0] {
			return x[i].lossVal[1] < x[j].lossVal[1]
		}
		return x[i].lossVal[0] < x[j].lossVal[0]
	})

	results := make([]int, len(x))
	for i := 0; i < len(x); i++ {
		results[i] = x[i].index
	}
	return results
}
