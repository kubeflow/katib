package random

import (
	"errors"
	"math/rand"

	"gonum.org/v1/gonum/floats"
)

// Multinomial draw samples from a multinomial distribution like numpy.random.multinomial.
// See https://docs.scipy.org/doc/numpy-1.15.0/reference/generated/numpy.random.multinomial.html
func Multinomial(n int, pvals []float64, size int) [][]int {
	result := make([][]int, size)
	l := len(pvals)
	x := make([]float64, l)
	floats.CumSum(x, pvals)

	for i := range result {
		result[i] = make([]int, l)

		for j := 0; j < n; j++ {

			var index int
			r := rand.Float64()
			for i := range x {
				if x[i] > r {
					index = i
					break
				}
			}
			result[i][index]++
		}
	}
	return result
}

// ArgMaxMultinomial returns the index sampled by multinomial distribution with given probabilities.
func ArgMaxMultinomial(pvals []float64) (int, error) {
	x := make([]float64, len(pvals))
	floats.CumSum(x, pvals)

	r := rand.Float64()
	for i := range x {
		if x[i] > r {
			return i, nil
		}
	}
	return 0, errors.New("invalid pvals")
}
