package main

import (
	"github.com/tuneinsight/lattigo/v5/core/rlwe"
	"github.com/tuneinsight/lattigo/v5/he/hefloat"
	"local.com/CaseStudy_Pack/src"
	"local.com/numerical-read-print/nread"
	"local.com/top-lattigo-ckks/pack"
)

func main() {

	disease := "breast"
	var num int
	if disease == "sepsis" {
		num = 138
	} else if disease == "breast"{
		num = 140
	} else {
		panic("disease name wrong.")
	}
	filename_data_matrix := "../../data/real/" + disease + "-test_data_x_y.txt"
	X := nread.M_Blank(filename_data_matrix)

	directory_WB := "../../model/"
	W_left, W_right, B_left, B_right := src.Read_WB_PDN(directory_WB, disease, "EPDN")

	var err error
	var params hefloat.Parameters
	if params, err = hefloat.NewParametersFromLiteral(
		hefloat.ParametersLiteral{
			LogN:            14,                                    // A ring degree of 2^{14}
			LogQ:            []int{55, 45, 45, 45, 45, 45, 45}, // An initial prime of 55 bits and 7 primes of 45 bits
			LogP:            []int{61},                             // The log2 size of the key-switching prime     // LogQP = 431  yangchen
			LogDefaultScale: 45,                                    // The default log2 of the scaling factor
		}); err != nil {
		panic(err)
	}

	kgen := rlwe.NewKeyGenerator(params)
	sk := kgen.GenSecretKeyNew()
	pk := kgen.GenPublicKeyNew(sk) 
	rlk := kgen.GenRelinearizationKeyNew(sk)
	evk := rlwe.NewMemEvaluationKeySet(rlk)

	encoder := hefloat.NewEncoder(params)
	encryptor := rlwe.NewEncryptor(params, pk)
	decryptor := rlwe.NewDecryptor(params, sk)
	eval := hefloat.NewEvaluator(params, evk)


	X_ct := pack.Encrypt_Matrix(params, encoder, encryptor, X)

	
	// left, err := pack.Mult_Matrix_Plain_Cipher(eval, W_left[0], X_ct)
	// if err != nil {
	// 	panic(err)
	// }
	// left = pack.Add_Vector(eval, B_left[0], left)

	// right, err := pack.Mult_Matrix_Plain_Cipher(eval, W_right[0], X_ct)
	// if err != nil {
	// 	panic(err)
	// }
	// right = pack.Add_Vector(eval, B_right[0], right)

	// result := pack.Mult(eval, left, right)


	result := X_ct   // 不会出错
	for i:=0; i<3; i++ {
		left, err := pack.Mult_Matrix_Plain_Cipher(eval, W_left[i], result)
		if err != nil {
			panic(err)
		}
		left = pack.Add_Vector(eval, B_left[i], left)
	
		right, err := pack.Mult_Matrix_Plain_Cipher(eval, W_right[i], result)
		if err != nil {
			panic(err)
		}
		right = pack.Add_Vector(eval, B_right[i], right)
	
		result = pack.Mult(eval, left, right)
	}

	pack.Print(decryptor, encoder, 6, num, result)
	
}