package pack

import (
	"github.com/tuneinsight/lattigo/v5/core/rlwe"
	"github.com/tuneinsight/lattigo/v5/he/hefloat"
)

func InnerProduct(eval *hefloat.Evaluator, coefficients []float64, input []*rlwe.Ciphertext) (output *rlwe.Ciphertext) {

	return Sum(eval, Mult_Vector(eval, coefficients, input))
}

func InnerProduct_AddBias(eval *hefloat.Evaluator, coefficients []float64, bias float64, input []*rlwe.Ciphertext) (output *rlwe.Ciphertext) {

	output = InnerProduct(eval, coefficients, input)
	if err := eval.Add(output, bias, output); err != nil {
		panic(err)
	}
	return output
}

func Add_Vector(eval *hefloat.Evaluator, vector []float64, ciphertexts_in []*rlwe.Ciphertext) (result []*rlwe.Ciphertext) {

	if len(vector) != len(ciphertexts_in) {
		panic("length of vector and ciphertexts_in not equal.")
	}

	result = make([]*rlwe.Ciphertext, len(ciphertexts_in))
	for i:=0;i<len(vector);i++ {
		var err error
		result[i], err = eval.AddNew(ciphertexts_in[i], vector[i])
		if err != nil {
			panic(err)
		}
	}
	return result
}

func Mult_Vector(eval *hefloat.Evaluator, vector []float64, ciphertexts_in []*rlwe.Ciphertext) (result []*rlwe.Ciphertext) {
	
	if len(vector) != len(ciphertexts_in) {
		panic("length of vector and ciphertexts_in not equal.")
	}
	
	result = make([]*rlwe.Ciphertext, len(ciphertexts_in))
	for i:=0;i<len(vector);i++ {
		var err error
		result[i], err = eval.MulRelinNew(ciphertexts_in[i], vector[i])
		if err != nil {
			panic(err)
		}
		if err = eval.Rescale(result[i], result[i]); err != nil {
			panic(err)
		}
	}
	return result
}