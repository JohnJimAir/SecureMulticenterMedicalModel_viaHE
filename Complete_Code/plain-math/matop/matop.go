package matop

import "local.com/plain-math/vecop"

func Multiply(mat0, mat1 [][]float64) (result [][]float64) {
	
	if ! IsMatrix(mat0) || ! IsMatrix(mat1) {
		panic("not matrix.")
	}

	if len(mat0[0]) != len(mat1) {
		panic("numCol of mat0 != numRow of mat1.")
	}
	
	numRow := len(mat0)
	numCol := len(mat1[0])
	result = make([][]float64, numRow)
	for i:=0; i<numRow; i++ {
		result[i] = make([]float64, numCol)
	}

	mat1_T := Transpose(mat1)
	for i:=0; i<numRow; i++ {
		for j:=0; j<numCol; j++ {
			result[i][j] = vecop.Dot(mat0[i], mat1_T[j])
		}
	}
	return result
}

func LinearTransform(mat [][]float64, vec []float64) (result []float64) {
	
	if !IsMatrix(mat) {
		panic("not matrix.")
	}
	if len(vec) == 0 {
		panic("length of vec = 0.")
	}
	if len(mat[0]) != len(vec) {
		panic("dimension not match.")
	}

	result = make([]float64, len(mat))
	for i:=0; i<len(mat); i++ {
		result[i] = vecop.Dot(mat[i], vec)
	}
	return result
} 

func Softmax_mat(mat [][]float64) (result [][]float64) {

	if !IsMatrix(mat) {
		panic("not matrix.")
	}
	
	for _, vec := range mat {
		result = append(result, vecop.Softmax(vec))
	}
	return result
}