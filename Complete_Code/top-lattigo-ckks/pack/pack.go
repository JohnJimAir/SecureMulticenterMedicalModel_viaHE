package pack

import (
	"github.com/tuneinsight/lattigo/v5/core/rlwe"
	"github.com/tuneinsight/lattigo/v5/he/hefloat"
)

func Add(eval *hefloat.Evaluator, ciphertexts_0 []*rlwe.Ciphertext, ciphertexts_1 []*rlwe.Ciphertext) (result []*rlwe.Ciphertext) {
	
	if len(ciphertexts_0) != len(ciphertexts_1) {
		panic("length of two slices of ciphertext not equal.")
	}

	length := len(ciphertexts_0)
	result = make([]*rlwe.Ciphertext, length)
	for i:=0;i<length;i++ {
		var err error
		result[i], err = eval.AddNew(ciphertexts_0[i], ciphertexts_1[i])
		if err != nil {
			panic(err)
		}
	}
	return result
}

func Sum(eval *hefloat.Evaluator, ciphertexts_in []*rlwe.Ciphertext) (output *rlwe.Ciphertext) {
	
	if len(ciphertexts_in) == 0 {
		panic("length of ciphertexts_in is 0")
	}

	output = (*ciphertexts_in[0]).CopyNew()
	for i:=1;i<len(ciphertexts_in);i++ {
		if err := eval.Add(output, ciphertexts_in[i], output); err != nil {
			panic(err)
		}
	}
	return output
}

func Mult(eval *hefloat.Evaluator, ciphertexts_0 []*rlwe.Ciphertext, ciphertexts_1 []*rlwe.Ciphertext) (result []*rlwe.Ciphertext) {
	
	if len(ciphertexts_0) != len(ciphertexts_1) {
		panic("length of two slices of ciphertext not equal.")
	}

	length := len(ciphertexts_0)
	result = make([]*rlwe.Ciphertext, length)
	for i:=0;i<length;i++ {
		var err error
		result[i], err = eval.MulRelinNew(ciphertexts_0[i], ciphertexts_1[i])
		if err != nil {
			panic(err)
		}
		if err = eval.Rescale(result[i], result[i]); err != nil {
			panic(err)
		}
	}
	return result
}