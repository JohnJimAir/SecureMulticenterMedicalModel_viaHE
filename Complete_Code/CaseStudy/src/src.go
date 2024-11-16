package src

import (
	"math/big"

	"github.com/tuneinsight/lattigo/v5/he"
	"github.com/tuneinsight/lattigo/v5/utils"
	"github.com/tuneinsight/lattigo/v5/utils/bignum"
)

// EvaluateLinearTransform evaluates a linear transform (i.e. matrix) on the input vector.
// values: the input vector
// diags: the non-zero diagonals of the linear transform
func EvaluateLinearTransform(values []complex128, diags map[int][]complex128) (res []complex128) {

	slots := len(values)
	keys := utils.GetKeys(diags)
	N1 := he.FindBestBSGSRatio(keys, len(values), 1)
	index, _, _ := he.BSGSIndex(keys, slots, N1)
	res = make([]complex128, slots)

	for j := range index {
		rot := -j & (slots - 1)
		tmp := make([]complex128, slots)

		for _, i := range index[j] {
			v, ok := diags[j+i]
			if !ok {
				v = diags[j+i-slots]
			}

			a := utils.RotateSlice(values, i)
			b := utils.RotateSlice(v, rot)

			for i := 0; i < slots; i++ {
				tmp[i] += a[i] * b[i]
			}
		}

		tmp = utils.RotateSlice(tmp, j)
		for i := 0; i < slots; i++ {
			res[i] += tmp[i]
		}
	}
	return
}

// GetChebyshevPoly returns the Chebyshev polynomial approximation of f the
// in the interval [-K, K] for the given degree.
func GetChebyshevPoly(K float64, degree int, f64 func(x float64) (y float64)) bignum.Polynomial {

	FBig := func(x *big.Float) (y *big.Float) {
		xF64, _ := x.Float64()
		return new(big.Float).SetPrec(x.Prec()).SetFloat64(f64(xF64))
	}

	var prec uint = 128

	interval := bignum.Interval{
		A:     *bignum.NewFloat(-K, prec),
		B:     *bignum.NewFloat(K, prec),
		Nodes: degree,
	}

	// Returns the polynomial.
	return bignum.ChebyshevApproximation(FBig, interval)
}