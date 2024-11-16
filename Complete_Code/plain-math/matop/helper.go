package matop

// Will return false, if the input is empty or dose not have a matrix format.
func IsMatrix(slice_2D [][]float64) (yn bool) {
	
	if len(slice_2D) == 0 || len(slice_2D[0]) == 0{
		return false
	}

	numColumn := len(slice_2D[0])
	for i:=1; i<len(slice_2D); i++ {
		if len(slice_2D[i]) != numColumn {
			return false
		}
	}
	return true
}

func Transpose(input [][]float64) (output [][]float64) {
	
	if IsMatrix(input) == false {
		panic("the input is not a matrix.")
	}

	numRow := len(input)
	numCol := len(input[0])
	output = make([][]float64, numCol)
	for i:=0;i<numCol;i++ {
		output[i] = make([]float64, numRow)
	}

	for i:=0;i<numCol;i++ {
		for j:=0;j<numRow;j++ {
			output[i][j] = input[j][i]
		}
	}
	return output
}