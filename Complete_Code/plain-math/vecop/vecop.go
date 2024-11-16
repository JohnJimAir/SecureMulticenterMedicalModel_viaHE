package vecop

func Mult(vec0, vec1 []float64) (result []float64) {

	if len(vec0) != len(vec1) || len(vec0) == 0 || len(vec1) == 0 {
		panic("length wrong.")
	}

	result = make([]float64, len(vec0))
	for i:=0; i<len(vec0); i++ {
		result[i] = vec0[i] * vec1[i]
	}
	return result
}

func Add(vec0, vec1 []float64) (result []float64) {

	if len(vec0) != len(vec1) || len(vec0) == 0 || len(vec1) == 0 {
		panic("length wrong.")
	}

	result = make([]float64, len(vec0))
	for i:=0; i<len(vec0); i++ {
		result[i] = vec0[i] + vec1[i]
	}
	return result
}

func Sub(vec0, vec1 []float64) (result []float64) {

	if len(vec0) != len(vec1) || len(vec0) == 0 || len(vec1) == 0 {
		panic("length wrong.")
	}

	result = make([]float64, len(vec0))
	for i:=0; i<len(vec0); i++ {
		result[i] = vec0[i] - vec1[i]
	}
	return result
}

func Dot(vec0, vec1 []float64) (result float64) {

	if len(vec0) != len(vec1) || len(vec0) == 0 || len(vec1) == 0 {
		panic("length wrong.")
	}

	result = 0
	for i:=0; i<len(vec0); i++ {
		result += vec0[i] * vec1[i]
	}
	return result
}