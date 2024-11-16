package pack

import (
	"math/big"

	"github.com/tuneinsight/lattigo/v5/core/rlwe"
	"github.com/tuneinsight/lattigo/v5/he/hefloat"
	"github.com/tuneinsight/lattigo/v5/utils/bignum"
)

// Need params to generate polyEval and get params.DefaultScale() for polyEval.Evaluate
func Polynomial(params hefloat.Parameters, eval *hefloat.Evaluator, f func(x float64) (y float64), K_left, K_right float64, degree int, ciphertexts_in []*rlwe.Ciphertext) (result []*rlwe.Ciphertext) {

	poly := hefloat.NewPolynomial(getChebyshevPoly(K_left, K_right, degree, f))
	polyEval := hefloat.NewPolynomialEvaluator(params, eval)
	scalar, constant := poly.ChangeOfBasis()

	result = make([]*rlwe.Ciphertext, len(ciphertexts_in))
	for i:=0; i<len(ciphertexts_in); i++ {
		var err error
		if err = eval.Mul(ciphertexts_in[i], scalar, ciphertexts_in[i]); err != nil {
			panic(err)
		}
		if err = eval.Add(ciphertexts_in[i], constant, ciphertexts_in[i]); err != nil {
			panic(err)
		}
		if err = eval.Rescale(ciphertexts_in[i], ciphertexts_in[i]); err != nil {
			panic(err)
		}

		if result[i], err = polyEval.Evaluate(ciphertexts_in[i], poly, params.DefaultScale()); err != nil {
			panic(err)
		}
	}
	return result
}

func getChebyshevPoly(K_left, K_right float64, degree int, f64 func(x float64) (y float64)) bignum.Polynomial {

	FBig := func(x *big.Float) (y *big.Float) {
		xF64, _ := x.Float64()
		return new(big.Float).SetPrec(x.Prec()).SetFloat64(f64(xF64))
	}

	var prec uint = 128
	interval := bignum.Interval{
		A:     *bignum.NewFloat(K_left, prec),
		B:     *bignum.NewFloat(K_right, prec),
		Nodes: degree,
	}

	return bignum.ChebyshevApproximation(FBig, interval)
}