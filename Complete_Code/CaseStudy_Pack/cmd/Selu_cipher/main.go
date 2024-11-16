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
	filename_data_matrix := "../../data/real/" + disease + "-test_data_x_y.txt"
	X := nread.M_Blank(filename_data_matrix)

	directory_WB := "../../model/"
	W, B := src.Read_WB_Selu(directory_WB, disease)

	var err error
	var params hefloat.Parameters
	if params, err = hefloat.NewParametersFromLiteral(
		hefloat.ParametersLiteral{
			LogN:            15,
			LogQ:            []int{51, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40},   // LogQP = 881 - 4*40 = 721  yangchen
			LogP:            []int{50, 50, 50},
			LogDefaultScale: 40,                                 // The default log2 of the scaling factor
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


	result := X_ct
	for i:=0; i<3; i++ {
		result, err = pack.Mult_Matrix_Plain_Cipher(eval, W[i], result)
		if err != nil {
			panic(err)
		}
		result = pack.Add_Vector(eval, B[i], result)
		if i<2 {
			result = pack.Polynomial(params, eval, src.Selu, -4.0, 4.0, 15, result)
		}
	}

	pack.Print(decryptor, encoder, 6, 140, result)

}

