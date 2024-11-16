package main

import (
	"local.com/CaseStudy_Pack/src"
	"local.com/numerical-read-print/nprint"
	"local.com/numerical-read-print/nread"
	"local.com/plain-math/matop"
	"local.com/plain-math/vecop"
)

func main() {

	disease := "sepsis"
	filename_data_matrix := "../../data/real/" + disease + "-test_data_x_y.txt"
	X := nread.M_Blank(filename_data_matrix)
	
	directory_WB := "../../model/"
	W_left, W_right, B_left, B_right := src.Read_WB_PDN(directory_WB, disease, "APDN")


	for i:=0; i<len(X); i++ {
		result := X[i]
		for layer:=0; layer<3; layer++ {
			left := matop.LinearTransform(W_left[layer], result)
			left = vecop.Add(B_left[layer], left)
		
			right := matop.LinearTransform(W_right[layer], X[i])
			right = vecop.Add(B_right[layer], right)
		
			result = vecop.Mult(left, right)
		}
		nprint.Print_Vector(result, 6, len(result))
	}
}