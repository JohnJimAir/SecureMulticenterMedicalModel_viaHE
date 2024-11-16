package src

import (
	"github.com/tuneinsight/lattigo/v5/core/rlwe"
	"github.com/tuneinsight/lattigo/v5/he/hefloat"
	"local.com/top-lattigo-ckks/pack"
)

type Node struct {
	Coefficients_mult []float64  // must not be integer
	Coefficient_add float64
	Activation func (float64) (float64)	
	Input []*rlwe.Ciphertext
}

func (n Node) Forward(interval []float64, degree int, eval *hefloat.Evaluator, params hefloat.Parameters) (output *rlwe.Ciphertext) {
	
	output = pack.InnerProduct_AddBias(eval, n.Coefficients_mult, n.Coefficient_add, n.Input)
	
	return pack.Polynomial(params, eval, n.Activation, interval[0], interval[1], degree, []*rlwe.Ciphertext{output})[0]
}
