package cmaes

import (
	"math"

	"gonum.org/v1/gonum/mat"
)

const minC = 1e-16

func initC(dim int) *mat.SymDense {
	c := mat.NewSymDense(dim, nil)
	for i := 0; i < dim; i++ {
		c.SetSym(i, i, 1.0)
	}
	return c
}

func floatsSqrtTo(src []float64) {
	for i := 0; i < len(src); i++ {
		src[i] = math.Sqrt(src[i])
	}
}

func sumColumns(m mat.Matrix) *mat.VecDense {
	r, c := m.Dims()
	x := make([]float64, r)
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			x[i] += m.At(i, j)
		}
	}
	return mat.NewVecDense(r, x)
}

func stackvec(r int, c int, vec *mat.VecDense) *mat.Dense {
	x := mat.NewDense(r, c, nil)
	for i := 0; i < r; i++ {
		x.SetRow(i, vec.RawVector().Data[:c])
	}
	return x
}

func arrinv(arr []float64) []float64 {
	x := make([]float64, len(arr))
	for i := 0; i < len(arr); i++ {
		x[i] = 1 / arr[i]
	}
	return x
}

func vecapply(vec *mat.VecDense, conv func(int, float64) float64) *mat.VecDense {
	x := mat.NewVecDense(vec.Len(), nil)
	for i := 0; i < vec.Len(); i++ {
		x.SetVec(i, conv(i, vec.AtVec(i)))
	}
	return x
}

func initMinC(dim int) *mat.SymDense {
	x := make([]float64, dim*dim)
	for i := 0; i < dim*dim; i++ {
		x[i] = minC
	}
	return mat.NewSymDense(dim, x)
}
