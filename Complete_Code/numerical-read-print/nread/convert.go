package nread

import "math"

func Int_from_float64(input [][]float64) (output [][]int) {
	
	output = make([][]int, len(input))

	for i:=0; i<len(input); i++ {
		output[i] = make([]int, len(input[i]))
		
		for j:=0; j<len(input[i]); j++ {
			output[i][j] = int(math.Round(input[i][j]))
		}
	}
	return output
}