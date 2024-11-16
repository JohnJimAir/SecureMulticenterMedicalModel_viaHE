package main

import (
	"fmt"
	"strconv"

	// "math"

	"github.com/tuneinsight/lattigo/v5/core/rlwe"
	"github.com/tuneinsight/lattigo/v5/he/hefloat"
	"local.com/CaseStudy/src"
	// "github.com/tuneinsight/lattigo/v5/utils/bignum"
)

func main() {
	// Reading data
	data_matrix_file_name := "../../data/3/sepsis-test_data_x_y_array1.txt"
	X, err1 := src.Readmatrix_nobracket(data_matrix_file_name)
	if err1 != nil {
		fmt.Println(err1)
		return
	}


	var W [][][]float64
	var B [][]float64
	for i:=1;i<=3;i++{// Layer i
		for j:=1; j<=2; j++{// The j-th linear transformation
			weight_matrix_file_name := "../../model/sepsis/APDN/W"+strconv.Itoa(i)+strconv.Itoa(j)+".txt"
			bias_vector_file_name := "../../model/sepsis/APDN/b"+strconv.Itoa(i)+strconv.Itoa(j)+".txt"
			w, err1 := src.Read_matrix(weight_matrix_file_name)
			if err1 != nil {
				fmt.Println(err1)
				return
			}

			W = append(W, w);

			b_tmp, err1 := src.Read_vector(bias_vector_file_name)
			if err1 != nil {
				fmt.Println(err1)
				return
			}
			// Print array
			// fmt.Println("The number "+ strconv.Itoa(j)+ " bias vector of layer " + strconv.Itoa(i) + ": ")
			// fmt.Println(b_tmp)

			B = append(B, b_tmp);
		}
	}

	// fmt.Println("The weight matrices: ")
	// for i, twoDSlice := range W {
    //     fmt.Printf("array[%d] = %v\n", i, twoDSlice)
    // }

	// fmt.Println("The bias vectors: ")
	// for i, twoDSlice := range B {
    //     fmt.Printf("array[%d] = %v\n", i, twoDSlice)
    // }



	// HE setup
	var err error
	var params hefloat.Parameters // Using the CKKS scheme
	if params, err = hefloat.NewParametersFromLiteral(
		hefloat.ParametersLiteral{
			// LogN:            14,                                    // A ring degree of 2^{14}
			// LogQ:            []int{55, 45, 45, 45, 45, 45, 45, 45}, // An initial prime of 55 bits and 7 primes of 45 bits
			// LogP:            []int{61},                             // The log2 size of the key-switching prime
			// LogDefaultScale: 45,                                    // The default log2 of the scaling factor
			LogN:            14,
			LogQ:            []int{45, 34, 34, 34, 34, 34, 34, 34, 34, 34},
			LogP:            []int{44, 43},
			LogDefaultScale: 34,
		}); err != nil {
		panic(err)
	}	
	// if params, err = hefloat.NewParametersFromLiteral(examples.HEFloatComplexParams[2]); err != nil {
	// 	panic(err)
	// }

	// prec := params.EncodingPrecision()

	kgen := rlwe.NewKeyGenerator(params)

	sk := kgen.GenSecretKeyNew()
	pk := kgen.GenPublicKeyNew(sk) 
	rlk := kgen.GenRelinearizationKeyNew(sk)
	evk := rlwe.NewMemEvaluationKeySet(rlk)

	LogSlots := params.LogMaxSlots()
	Slots := 1 << LogSlots // 2^LogSlots


	// Encoder
	ecd := hefloat.NewEncoder(params)
	// Encryptor
	enc := rlwe.NewEncryptor(params, pk)
	// Decryptor
	dec := rlwe.NewDecryptor(params, sk)
	// Evaluator
	eval := hefloat.NewEvaluator(params, evk)

	//To predict the k-th sample
	for k := 0; k < len(X); k++{
		// fmt.Println("\nThe "+strconv.Itoa(k) +"-th sample: ")
		// fmt.Println(X[k])
		x := src.Expand_vector(X[k], Slots)
		y := make([]float64, len(x))
    	copy(y, x)

		// Encoding
		pt_x := hefloat.NewPlaintext(params, params.MaxLevel())
		if err = ecd.Encode(x, pt_x); err != nil {
			panic(err)
		}
		pt_y := hefloat.NewPlaintext(params, params.MaxLevel())
		if err = ecd.Encode(y, pt_y); err != nil {
			panic(err)
		}

		//Encrypting
		ct_x, err := enc.EncryptNew(pt_x)
		if err != nil {
			panic(err)
		}

		ct_y, err := enc.EncryptNew(pt_y)
		if err != nil {
			panic(err)
		}

		ct_z, err := enc.EncryptNew(pt_x)
		if err != nil {
			panic(err)
		}
		// eval.Add(ct_y, ct_y, ct_y)
		// eval.Add(ct_z, ct_z, ct_z)

		for i:=1; i<=3; i++{// The i-th Layer
			// fmt.Println("\nLayer " + strconv.Itoa(i) + " begins ...")
			// the first linear transformation
			w := W[(i-1)*2];
			b := B[(i-1)*2];

			// Padding if the dimension is less than Slots
			n_rows := len(w)
			n_cols := len(w[0])
			w = src.Expand_matrix(w, Slots, Slots)
			b = src.Expand_vector(b, Slots)

			pt_b := hefloat.NewPlaintext(params, params.MaxLevel())
			if err = ecd.Encode(b, pt_b); err != nil {
				panic(err)
			} 

			size := n_cols-1 + n_rows
			nonZeroDiagonals := make([]int, size)

			for j := -n_rows + 1; j < n_cols; j++{
				nonZeroDiagonals[j+n_rows-1] = j
			}

			diagonals := make(hefloat.Diagonals[complex128])

			for _, l := range nonZeroDiagonals {
				tmp := make([]complex128, Slots)
				if l < 0 {
					for s := range tmp {
						if l+s+Slots < Slots {
							tmp[s] = complex(w[s][l+s+Slots],0)
						} else {
							tmp[s] = complex(w[s][(l+s+Slots)%Slots],0)
						}
					}
				}
				if l >= 0{
					for s := range tmp {
						if l+s < Slots {
							tmp[s] = complex(w[s][l+s],0)
						} else {
							tmp[s] = complex(w[s][(l+s)%Slots],0)
						}
					}
				}
				diagonals[l] = tmp
			}			

			ltparams := hefloat.LinearTransformationParameters{
				DiagonalsIndexList:       diagonals.DiagonalsIndexList(),
				Level:                    ct_x.Level(),
				Scale:                    rlwe.NewScale(params.Q()[ct_x.Level()]),
				LogDimensions:            ct_x.LogDimensions,
				LogBabyStepGianStepRatio: 1,
			}
		
			lt := hefloat.NewLinearTransformation(params, ltparams)
		
			if err := hefloat.EncodeLinearTransformation[complex128](ecd, diagonals, lt); err != nil {
				panic(err)
			}

			galEls := hefloat.GaloisElementsForLinearTransformation(params, ltparams)
		
			ltEval := hefloat.NewLinearTransformationEvaluator(eval.WithKey(rlwe.NewMemEvaluationKeySet(rlk, kgen.GenGaloisKeysNew(galEls, sk)...)))
		
			// ct_res := ct_x;
		
			if err := ltEval.Evaluate(ct_x, lt, ct_x); err != nil {
				panic(err)
			}
		
			// Result is not returned rescaled
			if err = eval.Rescale(ct_x, ct_x); err != nil {
				panic(err)
			}

			c_x := src.Float64_to_complex128(x);

			want := src.EvaluateLinearTransform(c_x, diagonals)
		
			// Decrypts the vector of plaintext values
			pt_res := dec.DecryptNew(ct_x)
		
			// Decodes the plaintext
			have := make([]float64, ct_x.Slots())
			if err = ecd.Decode(pt_res, have); err != nil {
				panic(err)
			}

			// Add the bias vector
			if err := eval.Add(ct_x, pt_b, ct_x); err != nil {
				panic(err)
			}

			// the second linear transformation
			w_1 := W[(i-1)*2+1];
			b_1 := B[(i-1)*2+1];

			// Padding if the dimension is less than Slots
			n_rows_1 := len(w_1)
			n_cols_1 := len(w_1[0])
			w_1 = src.Expand_matrix(w_1, Slots, Slots)
			b_1 = src.Expand_vector(b_1, Slots)
			
			pt_b_1 := hefloat.NewPlaintext(params, params.MaxLevel())
			if err = ecd.Encode(b_1, pt_b_1); err != nil {
				panic(err)
			} 
			
			size_1 := n_cols_1-1 + n_rows_1
			nonZeroDiagonals_1 := make([]int, size_1)
			
			for j := -n_rows_1 + 1; j < n_cols_1; j++{
				nonZeroDiagonals_1[j+n_rows_1-1] = j
			}
			// fmt.Println("The indices of the nonzero diagonals of the second weight matrix for layer "+ strconv.Itoa(i) + ": ")
			// fmt.Println(nonZeroDiagonals_1)
			
			diagonals_1 := make(hefloat.Diagonals[complex128])
			
			for _, l := range nonZeroDiagonals_1 {
				tmp := make([]complex128, Slots)
				if l < 0 {
					for s := range tmp {
						if l+s+Slots < Slots {
							tmp[s] = complex(w_1[s][l+s+Slots],0)
						} else {
							tmp[s] = complex(w_1[s][(l+s+Slots)%Slots],0)
						}
					}
				}
				if l >= 0{
					for s := range tmp {
						if l+s < Slots {
							tmp[s] = complex(w_1[s][l+s],0)
						} else {
							tmp[s] = complex(w_1[s][(l+s)%Slots],0)
						}
					}
				}
				diagonals_1[l] = tmp
			}			
			
			ltparams_1 := hefloat.LinearTransformationParameters{
				DiagonalsIndexList:       diagonals_1.DiagonalsIndexList(),
				Level:                    ct_z.Level(),
				Scale:                    rlwe.NewScale(params.Q()[ct_z.Level()]),
				LogDimensions:            ct_z.LogDimensions,
				LogBabyStepGianStepRatio: 1,
			}
			
			
			lt_1 := hefloat.NewLinearTransformation(params, ltparams_1)

			if err := hefloat.EncodeLinearTransformation[complex128](ecd, diagonals_1, lt_1); err != nil {
				panic(err)
			}	
			
			galEls_1 := hefloat.GaloisElementsForLinearTransformation(params, ltparams_1)
			
			ltEval_1 := hefloat.NewLinearTransformationEvaluator(eval.WithKey(rlwe.NewMemEvaluationKeySet(rlk, kgen.GenGaloisKeysNew(galEls_1, sk)...)))
			
			// ct_res := ct_x;
			
			if err := ltEval_1.Evaluate(ct_y, lt_1, ct_z); err != nil {
				panic(err)
			}
						
			// Result is not returned rescaled
			if err = eval.Rescale(ct_z, ct_z); err != nil {
				panic(err)
			}
					
			c_y := src.Float64_to_complex128(y);
			
			want_1 := src.EvaluateLinearTransform(c_y, diagonals_1)
			
			// Decrypts the vector of plaintext values
			pt_res_1 := dec.DecryptNew(ct_z)
			
			// Decodes the plaintext
			have_1 := make([]float64, ct_z.Slots())

			// Add the bias vector
			if err := eval.Add(ct_z, pt_b_1, ct_z); err != nil {
				panic(err)
			}

			pt_res_1 = dec.DecryptNew(ct_z)
			if err = ecd.Decode(pt_res_1, have_1); err != nil {
				panic(err)
			}
			


			
			ct_x, err = eval.MulRelinNew(ct_x, ct_z)
			if err != nil {
				panic(err)
			}

			if err := eval.Rescale(ct_x, ct_x); err != nil {
				panic(err)
			}

			for i:=0;i<len(x);i++ {
				x[i] = 0
			}
			for i := 0; i < n_rows && i < len(want); i++ {
				x[i] = real((want[i]+complex(b[i],0))*(want_1[i]+complex(b_1[i],0)))
				// fmt.Print((want[i]+complex(b[i],0))*(want_1[i]+complex(b_1[i],0)), ", ")
				// fmt.Print(x[i], ", ")
			}
			// fmt.Println()

			// Decodes the plaintext
			pt_res_2 := dec.DecryptNew(ct_x)
			have_2 := make([]float64, ct_x.Slots())
			if err = ecd.Decode(pt_res_2, have_2); err != nil {
				panic(err)
			}
			// fmt.Println("Encrypted results after multiplication in layer " + strconv.Itoa(i) + ": ")
			if i==3 {
				fmt.Println(have_2[:n_rows_1])
			}
			
		}
		// fmt.Println("********************************************************************")
	}
}
