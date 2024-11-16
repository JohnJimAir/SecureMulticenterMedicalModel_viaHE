package src

import (
	"math"
	"strconv"

	"local.com/numerical-read-print/nread"
)

func Selu(x float64) (y float64) {
	// return 1/(1+math.Exp(-x)) // sigmoid 
	// return (math.Exp(x)-math.Exp(-x))/(math.Exp(x)+math.Exp(-x)) // tanh
	if x < 0 { //SeLU
		return 1.05078 * 1.6733 * (math.Exp(x) - 1)
	} else {
		return 1.05078 * x
	}
}

func Selu_vector(x []float64) (y []float64) {
	if len(x) == 0 {
		panic("length = 0.")
	}
	y = make([]float64, len(x))
	for i:=0; i<len(x); i++ {
		y[i] = Selu(x[i])
	}
	return y
}

func Read_WB_PDN(directory, disease, PDN string) (W_left, W_right [][][]float64, B_left, B_right [][]float64) {

	for i:=1; i<=3; i++ { // Layer i
		for j:=1; j<=2; j++ { // The j-th linear transformation
			filename_weight := directory + disease + "/" + PDN + "/W"+ strconv.Itoa(i)+strconv.Itoa(j) +".txt"
			filename_bias := directory + disease + "/" + PDN + "/b"+ strconv.Itoa(i)+strconv.Itoa(j) +".txt"

			w := nread.M_Comma_Trim(filename_weight, ",[] ")
			if j == 1 {
				W_left = append(W_left, w)
			} else {
				W_right = append(W_right, w)
			}

			b := nread.V_Comma_Trim(filename_bias, ",[] ")
			if j == 1 {
				B_left = append(B_left, b)
			} else {
				B_right = append(B_right, b)
			}
		}
	}
	return 
}

func Read_WB_Selu(directory, disease string) (W [][][]float64, B [][]float64) {

	for i:=1; i<=3; i++ { // Layer i
		filename_weight := directory + disease + "/Selu/W"+ strconv.Itoa(i) +".txt"
		filename_bias := directory + disease + "/Selu/b"+ strconv.Itoa(i) +".txt"

		w := nread.M_Comma_Trim(filename_weight, ",[] ")
		W = append(W, w)

		b := nread.V_Comma_Trim(filename_bias, ",[] ")
		B = append(B, b)
	}
	return
}