package main

import (
	"fmt"
	"strconv"

	"math"

	"github.com/tuneinsight/lattigo/v5/core/rlwe"
	"github.com/tuneinsight/lattigo/v5/he/hefloat"
	"local.com/CaseStudy/src"
)


func main() {
	// Reading data
	data_matrix_file_name := "./data/syn3-test_data_x_y_2.txt"
	X, err1 := src.Readmatrix_nobracket(data_matrix_file_name)
	if err1 != nil {
		fmt.Println(err1)
		return
	}

	// Print 2D array
	// fmt.Println("Data matrix:")
	// for _, row := range X {
	// 	fmt.Println(row)
	// }

	var W [][][]float64
	for i := 1; i <= 3; i++ { // Layer i
		weight_matrix_file_name := "./model/Sigmoid/syn3/W" + strconv.Itoa(i) + ".txt"
		w, err1 := src.Read_matrix(weight_matrix_file_name)
		if err1 != nil {
			fmt.Println(err1)
			return
		}
		W = append(W, w)
	}

	// HE setup
	var err error
	var params hefloat.Parameters // Using the CKKS scheme
	if params, err = hefloat.NewParametersFromLiteral(
		hefloat.ParametersLiteral{
			LogN:            14,
			LogQ:            []int{42, 33, 33, 33, 33, 33, 33, 33, 33, 33, 33, 33}, // 42+12*33 = 438;
			LogP:            []int{33},
			LogDefaultScale: 33, // The default log2 of the scaling factor
		}); err != nil {
		panic(err)
	}

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
	for k := 0; k < len(X); k++ {
		// fmt.Println("\nThe " + strconv.Itoa(k) + "-th sample: ")
		// fmt.Println(X[k])
		x := src.Expand_vector(X[k], Slots)

		// Encoding
		pt_x := hefloat.NewPlaintext(params, params.MaxLevel())
		if err = ecd.Encode(x, pt_x); err != nil {
			panic(err)
		}

		//Encrypting
		ct_x, err := enc.EncryptNew(pt_x)
		if err != nil {
			panic(err)
		}

		// eval.Add(ct_y, ct_y, ct_y)
		// eval.Add(ct_z, ct_z, ct_z)

		for i := 1; i <= 3; i++ { // The i-th Layer
			// fmt.Println("\nLayer " + strconv.Itoa(i) + " begins ...")
			// the first linear transformation
			w := W[i-1]
			// fmt.Println("The first weight matrix of layer " + strconv.Itoa(i) + ": ")
			// for _, row := range w {
			// 	fmt.Println(row)
			// }

			// Padding if the dimension is less than Slots
			n_rows := len(w)
			n_cols := len(w[0])
			w = src.Expand_matrix(w, Slots, Slots)

			size := n_cols - 1 + n_rows
			nonZeroDiagonals := make([]int, size)

			for j := -n_rows + 1; j < n_cols; j++ {
				nonZeroDiagonals[j+n_rows-1] = j
			}
			// fmt.Println("The indices of the nonzero diagonals of the first weight matrix for layer " + strconv.Itoa(i) + ": ")
			// fmt.Println(nonZeroDiagonals)

			diagonals := make(hefloat.Diagonals[complex128])

			for _, l := range nonZeroDiagonals {
				tmp := make([]complex128, Slots)
				if l < 0 {
					for s := range tmp {
						if l+s+Slots < Slots {
							tmp[s] = complex(w[s][l+s+Slots], 0)
						} else {
							tmp[s] = complex(w[s][(l+s+Slots)%Slots], 0)
						}
					}
				}
				if l >= 0 {
					for s := range tmp {
						if l+s < Slots {
							tmp[s] = complex(w[s][l+s], 0)
						} else {
							tmp[s] = complex(w[s][(l+s)%Slots], 0)
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

			c_x := src.Float64_to_complex128(x)

			want := src.EvaluateLinearTransform(c_x, diagonals)
			x = src.Real_part_of_complex_array(want)

			// fmt.Println("Plaintext results after the first linear tranformation of layer " + strconv.Itoa(i) + ": ")
			if i==3 {
				for i := 0; i < n_rows && i < len(want); i++ {
					fmt.Print(x[i], ", ")
				}
				fmt.Println()
			}


			// Decrypts the vector of plaintext values
			pt_res := dec.DecryptNew(ct_x)

			// Decodes the plaintext
			have := make([]float64, ct_x.Slots())
			if err = ecd.Decode(pt_res, have); err != nil {
				panic(err)
			}

			// fmt.Println("Encrypted results after the first linear tranformation of layer " + strconv.Itoa(i) + ": ")
			// fmt.Println(have[:n_rows])

			// fmt.Printf("vector x matrix %s", hefloat.GetPrecisionStats(params, ecd, dec, want, ct_x, 0, false).String())

			if i < 3 {
				// Activation

				sigmoid := func(x float64) (y float64) {
					return 1/(1+math.Exp(-x)) // sigmoid 
					// return (math.Exp(x)-math.Exp(-x))/(math.Exp(x)+math.Exp(-x)) // tanh
					// if x < 0 { //SeLU
					// 	return 1.05078 * 1.6733 * (math.Exp(x) - 1)
					// } else {
					// 	return 1.05078 * x
					// }
				}

				K := 2.0
				// Chebyhsev approximation of the activation in the domain [-K, K] of degree 3.
				poly := hefloat.NewPolynomial(src.GetChebyshevPoly(K, 7, sigmoid))
				// poly := hefloat.NewPolynomial(GetChebyshevPoly(K, 63, sigmoid))

				// Instantiates the polynomial evaluator
				polyEval := hefloat.NewPolynomialEvaluator(params, eval)

				// Retrieves the change of basis y = scalar * x + constant
				scalar, constant := poly.ChangeOfBasis()

				// Performes the change of basis Standard -> Chebyshev
				if err := eval.Mul(ct_x, scalar, ct_x); err != nil {
					panic(err)
				}

				if err := eval.Add(ct_x, constant, ct_x); err != nil {
					panic(err)
				}

				if err := eval.Rescale(ct_x, ct_x); err != nil {
					panic(err)
				}

				// Evaluates the polynomial
				if ct_x, err = polyEval.Evaluate(ct_x, poly, params.DefaultScale()); err != nil {
					panic(err)
				}

				want_activation := make([]float64, ct_x.Slots())
				want_sigmoid := make([]float64, ct_x.Slots())
				for i := 0; i < n_rows && i < len(want_activation); i++ {
					want_activation[i], _ = poly.Evaluate(x[i])[0].Float64()
					want_sigmoid[i] = sigmoid(x[i])
				}
				// x = want_activation
				x = want_sigmoid
				// fmt.Println("\n Plaintext results after activation in layer " + strconv.Itoa(i) + ": ")
				// fmt.Println(want_activation[:n_rows])
				// fmt.Println("\n Plaintext results after sigmoid in layer " + strconv.Itoa(i) + ": ")
				// fmt.Println(want_sigmoid[:n_rows])

				// Decrypt and Decodes
				pt_res = dec.DecryptNew(ct_x)
				have_activation := make([]float64, ct_x.Slots())
				if err = ecd.Decode(pt_res, have_activation); err != nil {
					panic(err)
				}

				// fmt.Println("Encrypted results after activation in layer " + strconv.Itoa(i) + ": ")
				// fmt.Println(have_activation[:n_rows])
			}
		}

	}
	// fmt.Println("********************************************************************")
}
