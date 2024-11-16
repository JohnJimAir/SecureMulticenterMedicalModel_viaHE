package nprint

import "fmt"

func Print_Matrix(mat [][]float64, precision int, num int) {
	
	if num > len(mat) {
		num = len(mat)
	}
	for i:=0; i<num; i++ {
		format := fmt.Sprintf("%%.%df,", precision)
		for j:=0; j<len(mat[i]); j++ {
			if j != len(mat[i])-1 {
				fmt.Printf(format, mat[i][j])
			} else {
				fmt.Printf(format[:len(format)-1], mat[i][j])
			}
		}
		fmt.Println()
	}
}

func Print_Vector(vec []float64, precision int, num int) {
	
	if num > len(vec) {
		num = len(vec)
	}
	for i:=0; i<num; i++ {
		format := fmt.Sprintf("%%.%df,", precision)
		if i != num-1 {
			fmt.Printf(format, vec[i])
		} else {
			fmt.Printf(format[:len(format)-1], vec[i])
		}
	}
	fmt.Println()
}