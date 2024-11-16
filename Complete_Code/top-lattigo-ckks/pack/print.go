package pack

import (
	"github.com/tuneinsight/lattigo/v5/core/rlwe"
	"github.com/tuneinsight/lattigo/v5/he/hefloat"
	"local.com/numerical-read-print/nprint"
)

func Print(decryptor *rlwe.Decryptor, encoder *hefloat.Encoder, precision int, num int, ciphertexts_in []*rlwe.Ciphertext) {
	
	result := Decrypt_Matrix(decryptor, encoder, num, ciphertexts_in)

	nprint.Print_Matrix(result, precision, num)
}