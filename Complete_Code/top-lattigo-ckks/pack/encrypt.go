package pack

import (
	"github.com/tuneinsight/lattigo/v5/core/rlwe"
	"github.com/tuneinsight/lattigo/v5/he/hefloat"
	"local.com/plain-math/matop"
)

// One ciphertext corresponds to one column.
func Encrypt_Matrix(params hefloat.Parameters, encoder *hefloat.Encoder, encryptor *rlwe.Encryptor, matrix [][]float64) (output []*rlwe.Ciphertext) {
	
	if len(matrix) > params.MaxSlots() {
		panic("exceed maxslots.")
	}

	matrix_T := matop.Transpose(matrix)
	length := len(matrix_T)

	output = make([]*rlwe.Ciphertext, length)
	pt := hefloat.NewPlaintext(params, params.MaxLevel())
	var err error
	for i:=0;i<length;i++ {
		if err = encoder.Encode(matrix_T[i], pt); err != nil {
			panic(err)
		}
		output[i], err = encryptor.EncryptNew(pt)
		if err != nil {
			panic(err)
		}
	}
	return output
}

func Decrypt_Matrix(decryptor *rlwe.Decryptor, encoder *hefloat.Encoder, num int, ciphertexts_in []*rlwe.Ciphertext) ([][]float64) {

	length := len(ciphertexts_in)

	matrix := make([][]float64, length)
	for i:=0; i<length; i++ {	
		pt := decryptor.DecryptNew(ciphertexts_in[i])

		values := make([]float64, ciphertexts_in[i].Slots())
		if err := encoder.Decode(pt, values); err != nil {
			panic(err)
		}
		matrix[i] = values[0:num]
	}
	return matop.Transpose(matrix)
}

func Decrypt_Matrix_full(decryptor *rlwe.Decryptor, encoder *hefloat.Encoder, ciphertexts_in []*rlwe.Ciphertext) ([][]float64) {
	
	return Decrypt_Matrix(decryptor, encoder, ciphertexts_in[0].Slots(), ciphertexts_in)
}