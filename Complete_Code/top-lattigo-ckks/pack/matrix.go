package pack

import (
	"errors"

	"github.com/tuneinsight/lattigo/v5/core/rlwe"
	"github.com/tuneinsight/lattigo/v5/he/hefloat"
	"local.com/plain-math/matop"
)

func Mult_Matrix_Plain_Cipher(eval *hefloat.Evaluator, plainMatrix [][]float64, ciphertexts_in []*rlwe.Ciphertext) (result []*rlwe.Ciphertext, err error) {
	
	if matop.IsMatrix(plainMatrix) == false {
		return nil, errors.New("The plainMatrix is not a matrix.")
	}
	if len(plainMatrix[0]) != len(ciphertexts_in) {
		return nil, errors.New("Dimension not match.")
	}

	numRow_plainMatrix := len(plainMatrix)
	result = make([]*rlwe.Ciphertext, numRow_plainMatrix)
	for i:=0;i<numRow_plainMatrix;i++ {
		result[i] = InnerProduct(eval, plainMatrix[i], ciphertexts_in)
	}
	return result, nil
}

func Mult_Matrix_Cipher_Plain(eval *hefloat.Evaluator, ciphertexts_in []*rlwe.Ciphertext, plainMatrix [][]float64) (result []*rlwe.Ciphertext, err error) {
	
	if matop.IsMatrix(plainMatrix) == false {
		return nil, errors.New("The plainMatrix is not a matrix.")
	}
	if len(plainMatrix[0]) != len(ciphertexts_in) {
		return nil, errors.New("Dimension not match.")
	}

	plainMatrix_T := matop.Transpose(plainMatrix)
	numCol_plainMatrix := len(plainMatrix[0])
	result = make([]*rlwe.Ciphertext, numCol_plainMatrix)
	for i:=0; i<numCol_plainMatrix; i++ {
		result[i] = InnerProduct(eval, plainMatrix_T[i], ciphertexts_in)
	}
	return result, nil
}