// Package main implements an example of smooth function approximation using Chebyshev polynomial interpolation.
package main

import (
	"math"

	"github.com/JohnJimAir/kan-he/src"
	"github.com/tuneinsight/lattigo/v5/core/rlwe"
	"github.com/tuneinsight/lattigo/v5/he/hefloat"
	"local.com/numerical-read-print/nread"
	"local.com/top-lattigo-ckks/pack"
)

func main() {

	var err error
	filename := "../../data/test_data_breast-cancer.csv"
	input := nread.M_Comma_EmitFirstline(filename)
	

	params, err := hefloat.NewParametersFromLiteral(hefloat.ParametersLiteral{
		LogN:            16,
		LogQ:            []int{56, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45,}, // logQP = 1356, yangchen
		LogP:            []int{55, 55, 55, 55},
		LogDefaultScale: 45,
	})
	if err != nil {
		panic(err)
	}

	kgen := rlwe.NewKeyGenerator(params)
	sk, pk := kgen.GenKeyPairNew()
	encoder := hefloat.NewEncoder(params)
	decryptor := rlwe.NewDecryptor(params, sk)
	encryptor := rlwe.NewEncryptor(params, pk)

	rlk := kgen.GenRelinearizationKeyNew(sk)
	evk := rlwe.NewMemEvaluationKeySet(rlk)
	eval := hefloat.NewEvaluator(params, evk)


	tanh := func(x float64) (y float64) { return math.Tanh(x) }
	sin := func(x float64) (y float64) { return math.Sin(x) }
	abs := func(x float64) (y float64) { return math.Abs(x) }
	exp := func(x float64) (y float64) { return math.Exp(x) }
	pow2 := func(x float64) (y float64) { return math.Pow(x, 2) }
	pow3 := func(x float64) (y float64) {return math.Pow(x, 3) }
	identity := func(x float64) (y float64) { return x }
	contract := func(x float64) (y float64) { return 0.000001* x }
	// log := func(x float64) (y float64) { return math.Log(x) }
	// sqrt := func(x float64) (y float64) { return math.Sqrt(x) }
	


	K := 16.0

	input_ct := pack.Encrypt_Matrix(params, encoder, encryptor, input)
	input_ct_2d := make([][]*rlwe.Ciphertext, 0)
	for i:=0;i<len(input_ct);i++ {
		input_ct_2d = append(input_ct_2d, []*rlwe.Ciphertext{input_ct[i]})
	}
	
	var blo_top_0 *src.Block = new(src.Block)
	blo_top_0.Initialize(9, [][]float64{{3.77}, {7.07}, {9.52}, {9.96}, {3.64}, {2.24}, {10.000001}, {7.85}, {7.94}}, 
		[]float64{-1.01, -6.21, -8.15, -3.26, -0.62, 8.2, -8.2, 7.58, -0.2}, 
		[]func (float64) (float64){tanh, sin, sin, abs, sin, sin, tanh, sin, abs},
		input_ct_2d,
	)
	out_blo_top_0 := blo_top_0.Forward(src.Interval([]float64{K,K,K, 8,K,K, K,K,8}), []int{31,31,31, 31,31,31, 31,31,31}, eval, params)

	var blo_top_1 *src.Block = new(src.Block)
	blo_top_1.Initialize(9, [][]float64{{0.000001}, {7.4}, {-1.000001}, {9.6}, {6.44}, {6.11}, {5.2}, {4.95}, {5.89}}, 
		[]float64{0.000001, 1.19, 0.33, -2.47, -2.23, -0.73, 1.18, 9.62, -2.45}, 
		[]func (float64) (float64){contract, sin, pow2, tanh, sin, sin, sin, sin, tanh},
		input_ct_2d,
	)
	out_blo_top_1 := blo_top_1.Forward(src.Interval([]float64{K,K,K, K,K,K, K,K,K}), []int{31,31,31, 31,31,31, 31,31,31}, eval, params)

	var blo_top_2 *src.Block = new(src.Block)
	blo_top_2.Initialize(9, [][]float64{{-1.000001}, {5.08}, {6.62}, {7.21}, {2.2}, {3.24}, {-1.000001}, {-1.000001}, {0.28}}, 
		[]float64{0.43, -2.22, 2.99, -5.79, -9.64, -2.6, 0.24, 0.37, 1.0}, 
		[]func (float64) (float64){pow3, sin, sin, sin, contract, tanh, pow2, pow3, contract},
		input_ct_2d,
	)
	out_blo_top_2 := blo_top_2.Forward(src.Interval([]float64{K,K,K, K,K,K, K,K,K}), []int{31,31,31, 31,31,31, 31,31,31}, eval, params)

	var blo_top_3 *src.Block = new(src.Block)
	blo_top_3.Initialize(9, [][]float64{{1.13}, {3.94}, {3.89}, {3.86}, {1.49}, {10.000001}, {3.65}, {9.79}, {7.8}}, 
		[]float64{-9.75, -0.58, -7.86, -8.02, 2.53, -2.6, -1.43, 4.21, -0.84}, 
		[]func (float64) (float64){contract, tanh, sin, sin, contract, tanh, sin, sin, tanh},
		input_ct_2d,
	)
	out_blo_top_3 := blo_top_3.Forward(src.Interval([]float64{K,K,K, K,K,K, K,K,K}), []int{31,31,31, 31,31,31, 31,31,31}, eval, params)


	var blo_middle *src.Block = new(src.Block)
	coefficient_middle := [][]float64{
		{0.08, 0.39, 0.09, 0.000001/*-0.e-2*/, 0.21, -0.13, 0.19, 0.01, 0.01},
		{0.000001, 1.38, 2.28, 0.27, 1.64, 0.72, -0.37, -0.87, 0.29},
		{-0.61, -0.01, -0.04, 0.05, 0.000001/*tan -0.e-2*/, 0.07, 0.05, -0.3, 0.000001/*tan*/},
		{0.000001/*tan*/, 25.59, 21.1, 19.17, 0.000001/*tan*/, 12.94, 4.29, 12.29, 9.55},
	}
	blo_middle.Initialize(4, coefficient_middle, 
		[]float64{-1.0, 5.55, 4.04, 85.59}, 
		[]func (float64) (float64){pow2, sin, sin, identity},
		[][]*rlwe.Ciphertext{out_blo_top_0, out_blo_top_1, out_blo_top_2, out_blo_top_3},
	)
	out_blo_middle := blo_middle.Forward(src.Interval([]float64{K,K,K, K}), []int{31,31,31, 31}, eval, params)


	var blo_bottom *src.Block = new(src.Block)
	coefficient_bottom := [][]float64 {
		{0.13, -0.09, 2.93, -0.01},
		{0.31, -0.21, 7.04, -0.03},
	}
	blo_bottom.Initialize(2, coefficient_bottom, 
		[]float64{0.0, 10.72}, 
		[]func (float64) (float64){exp, tanh},
		[][]*rlwe.Ciphertext{out_blo_middle, out_blo_middle},
	)
	out_blo_bottom := blo_bottom.Forward(src.Interval([]float64{8,8}), []int{31,31}, eval, params)


	var blo_trick *src.Block = new(src.Block)
	blo_trick.Initialize(2, [][]float64{{1988.48}, {-7.34}}, 
		[]float64{-31.97, 1.99}, 
		[]func (float64) (float64){identity, identity},
		[][]*rlwe.Ciphertext{ {out_blo_bottom[0]}, {out_blo_bottom[1]} },
	)
	out_blo_trick := blo_trick.Forward(src.Interval([]float64{8,8}), []int{1,1}, eval, params)

	pack.Print(decryptor, encoder, 6, 140, out_blo_trick)

}