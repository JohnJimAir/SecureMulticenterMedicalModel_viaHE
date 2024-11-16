// Package main implements an example of smooth function approximation using Chebyshev polynomial interpolation.
package main

import (
	"fmt"
	"math"
	"math/big"

	"github.com/JohnJimAir/kan-he/src"
	"github.com/tuneinsight/lattigo/v5/core/rlwe"
	"github.com/tuneinsight/lattigo/v5/he/hefloat"
	"github.com/tuneinsight/lattigo/v5/utils/bignum"
	"local.com/numerical-read-print/nread"
	"local.com/top-lattigo-ckks/pack"
)

func main() {
	
	var err error
	filename := "../../data/test_data_sepsis.csv"
	input := nread.M_Comma_EmitFirstline(filename)
	

	params, err := hefloat.NewParametersFromLiteral(hefloat.ParametersLiteral{
		// LogN:            LogN,                                              // Log2 of the ring degree
		// LogQ:            []int{55, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40}, // Log2 of the ciphertext prime moduli
		// LogP:            []int{61, 61, 61},                                 // Log2 of the key-switch auxiliary prime moduli
		// LogDefaultScale: 40,                                                // Log2 of the scale
		// Xs:              ring.Ternary{H: 192},
		LogN:            15,
		LogQ:            []int{51, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40},
		LogP:            []int{50, 50, 50},
		LogDefaultScale: 40,
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
	// abs := func(x float64) (y float64) { return math.Abs(x) }
	// exp := func(x float64) (y float64) { return math.Exp(x) }
	pow2 := func(x float64) (y float64) { return math.Pow(x, 2) }
	// pow3 := func(x float64) (y float64) {return math.Pow(x, 3) }
	identity := func(x float64) (y float64) { return x }
	contract := func(x float64) (y float64) { return 0.000001* x }
	log := func(x float64) (y float64) { return math.Log(x) }
	sqrt := func(x float64) (y float64) { return math.Sqrt(x) }


	// An easy example to test the correctness of Node.Forward
	pt := hefloat.NewPlaintext(params, params.MaxLevel())
	values := make([]float64, pt.Slots())
	for i := range values {
		values[i] = 2.3
	}
	if err = encoder.Encode(values, pt); err != nil {
		panic(err)
	}
	var ct_in *rlwe.Ciphertext
	if ct_in, err = encryptor.EncryptNew(pt); err != nil {
		panic(err)
	}

	nn := src.Node{
		Coefficients_mult: []float64{1.000001},
		Coefficient_add: 0.000001,
		Activation : tanh,
		Input: []*rlwe.Ciphertext{ct_in},
	}
	out := nn.Forward([]float64{-4.0, 4.0}, 31, eval, params)
	fmt.Println("Using Node.Forward: ")
	pack.Print(decryptor, encoder, 8, 4, []*rlwe.Ciphertext{out})


	poly := hefloat.NewPolynomial(GetChebyshevPoly(4, 31, tanh))
	polyEval := hefloat.NewPolynomialEvaluator(params, eval)
	scalar, constant := poly.ChangeOfBasis()

	var ct_out *rlwe.Ciphertext
	if ct_out, err = encryptor.EncryptNew(pt); err != nil {
		panic(err)
	}
	if err := eval.Mul(ct_in, scalar, ct_out); err != nil {
		panic(err)
	}
	if err := eval.Add(ct_out, constant, ct_out); err != nil {
		panic(err)
	}
	if err := eval.Rescale(ct_out, ct_out); err != nil {
		panic(err)
	}

	if ct_out, err = polyEval.Evaluate(ct_out, poly, params.DefaultScale()); err != nil {
		panic(err)
	}
	want := make([]float64, ct_out.Slots())
	for i := range want {
		// want[i], _ = poly.Evaluate(values[i])[0].Float64()
		want[i] = tanh(values[i])
	}
	fmt.Println("The difference from polynomial approximation: ")
	PrintPrecisionStats(params, ct_out, want, encoder, decryptor)

	// The easy example ends



	K := 16.0

	input_ct := pack.Encrypt_Matrix(params, encoder, encryptor, input)
	input_ct_2d := make([][]*rlwe.Ciphertext, 0)
	for i:=0;i<len(input_ct);i++ {
		input_ct_2d = append(input_ct_2d, []*rlwe.Ciphertext{input_ct[i]})
	}


	var blo_top *src.Block = new(src.Block)
	coefficient_top_mult := [][]float64{{-1.0000001},{0.58},{0.4},{0.14},{0.2},
		{0.0000001},{0.0000001},{1.06},{0.0000001},{0.28},
		{1.0000001},{0.0000001},{0.0000001},{3.4},{0.0000001},
		{-1.38},{0.27},{9.96},{0.31},{0.89},
		{0.95},{0.43},{0.41},{0.31},{0.16},
		{0.18},{0.17},{1.02},{0.21},{0.23},
		{0.35},{1.22},{1.51},{6.07},{0.26},
		{0.28},{0.26}}
	coefficient_top_add := []float64{-0.72, -0.48, 1.37, 1.0, -0.85,
		0.0, 0.0, -9.61, 0.0, -5.95, 
		0.37, 0.0, 0.0, 3.95, 0.0,
		4.25, 1.85, 7.21, 5.04, -0.18, 
		-0.53, 2.24, 2.39, 1.61, -4.2,
		-7.56, 8.5, -0.68, 2.16, -7.04,
		-1.43, -2.35, -2.12, 2.42, 2.15,
		-7.78, 4.52}
	blo_top.Initialize(37, coefficient_top_mult, coefficient_top_add, 
		[]func (float64) (float64){pow2, tanh, sin, contract, tanh,
			contract, contract, sin, contract, contract,
			sqrt, contract, contract, log, contract,
			log, sin, identity, sin, sin,
			tanh, sin, sin, sin, sin,
			sin, sin, tanh, sin,sin,
			tanh, tanh, tanh, identity, sin,
			sin, sin},
		input_ct_2d )
	out_blo_top := blo_top.Forward([][]float64{{-K,K},{-K,K},{-K,K},{-K,K},{-K,K}, {-K,K},{-K,K},{-K,K},{-K,K},{-K,K}, {0.0,K},{-K,K},{-K,K},{0.0,K},{-K,K},
		{0.0,K},{-K,K},{-K,K},{-K,K},{-K,K}, {-K,K},{-K,K},{-K,K},{-K,K},{-K,K}, {-K,K},{-K,K},{-K,K},{-K,K},{-K,K}, {-K,K},{-K,K},{-K,K},{-K,K},{-K,K}, {-K,K},{-K,K}}, 
		[]int{31,31,31,31,31, 31,31,31,31,31, 31,31,31,31,31, 31,31,31,31,31, 31,31,31,31,31, 31,31,31,31,31, 31,31,31,31,31, 31,31}, 
		eval, params)


	var blo_bottom *src.Block = new(src.Block)
	coefficient_bottom_mult := []float64{0.04, -0.07, -0.05, 0.01, -0.32, 
		0.0, 0.0, 0.05, 0.0, 0.1, 
		-0.28, 0.0, 0.0, -0.24, 0.0,
		0.05, 0.12, 0.02, -0.12, 0.02,
		0.03, -0.49, 0.13, 0.24, 0.35,
		0.13, 0.5, 0.02, 0.13, -0.21,  
		-0.11, 0.04, 0.03, -0.01, 0.18,
		-0.06, 0.04}
	blo_bottom.Initialize(1, [][]float64{coefficient_bottom_mult}, 
		[]float64{5.76}, 
		[]func (float64) (float64){sin},
		[][]*rlwe.Ciphertext{out_blo_top})
	out_blo_bottom := blo_bottom.Forward([][]float64{{K,-K}}, []int{31}, eval, params)


	var blo_trick *src.Block = new(src.Block)
	blo_trick.Initialize(1, [][]float64{{-1.05}}, 
		[]float64{1.04}, 
		[]func (float64) (float64){identity},
		[][]*rlwe.Ciphertext{ {out_blo_bottom[0]} },
	)
	out_blo_trick := blo_trick.Forward([][]float64{{-K,K}}, []int{1}, eval, params)


	pack.Print(decryptor, encoder, 8, 138, out_blo_trick)

}

// GetChebyshevPoly returns the Chebyshev polynomial approximation of f the
// in the interval [-K, K] for the given degree.
func GetChebyshevPoly(K float64, degree int, f64 func(x float64) (y float64)) bignum.Polynomial {

	FBig := func(x *big.Float) (y *big.Float) {
		xF64, _ := x.Float64()
		return new(big.Float).SetPrec(x.Prec()).SetFloat64(f64(xF64))
	}

	var prec uint = 128

	interval := bignum.Interval{
		A:     *bignum.NewFloat(-K, prec),
		B:     *bignum.NewFloat(K, prec),
		Nodes: degree,
	}

	// Returns the polynomial.
	return bignum.ChebyshevApproximation(FBig, interval)
}

// PrintPrecisionStats decrypts, decodes and prints the precision stats of a ciphertext.
func PrintPrecisionStats(params hefloat.Parameters, ct *rlwe.Ciphertext, want []float64, ecd *hefloat.Encoder, dec *rlwe.Decryptor) {

	var err error

	// Decrypts the vector of plaintext values
	pt := dec.DecryptNew(ct)

	// Decodes the plaintext
	have := make([]float64, ct.Slots())
	if err = ecd.Decode(pt, have); err != nil {
		panic(err)
	}

	// Pretty prints some values
	fmt.Printf("Have: ")
	for i := 0; i < 4; i++ {
		fmt.Printf("%20.15f ", have[i])
	}
	fmt.Printf("...\n")

	fmt.Printf("Want: ")
	for i := 0; i < 4; i++ {
		fmt.Printf("%20.15f ", want[i])
	}
	fmt.Printf("...\n")

	// Pretty prints the precision stats
	// fmt.Println(hefloat.GetPrecisionStats(params, ecd, dec, have, want, 0, false).String())
}

