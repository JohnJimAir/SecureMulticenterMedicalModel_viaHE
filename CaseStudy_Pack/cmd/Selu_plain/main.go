package main

import (
	"local.com/CaseStudy_Pack/src"
	"local.com/numerical-read-print/nprint"
	"local.com/numerical-read-print/nread"
	"local.com/plain-math/matop"
	"local.com/plain-math/vecop"
)

func main() {
	
	disease := "breast"
	filename_data_matrix := "../../data/real/" + disease + "-test_data_x_y.txt"
	X := nread.M_Blank(filename_data_matrix)

	directory_WB := "../../model/"
	W, B := src.Read_WB_Selu(directory_WB, disease)

	for i:=0; i<len(X); i++ {
		result := X[i]
		for layer:=0; layer<3; layer++ {
			result = matop.LinearTransform(W[layer], result)
			result = vecop.Add(B[layer], result)
			if i<2 {
				result = src.Selu_vector(result)
			}
		}
		nprint.Print_Vector(result, 6, len(result))
	}
}