package src

// expandMatrix expands a 2D array A to a p*q 2D array B, with all new elements set to 0
func Expand_matrix(A [][]float64, p int, q int) [][]float64 {
	m := len(A)
	n := len(A[0])

	// Initialize p*q 2D array B with all elements set to 0
	B := make([][]float64, p)
	for i := range B {
		B[i] = make([]float64, q)
	}

	// Copy elements from A to the first m rows and n columns of B
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			B[i][j] = A[i][j]
		}
	}
	return B
}

// expandArray expands a 1D array a to a length n 1D array b, with all new elements set to 0
func Expand_vector(a []float64, n int) []float64 {
	m := len(a)

	// Initialize array b with length n and all elements set to 0
	b := make([]float64, n)

	// Copy elements from a to the first m positions of b
	for i := 0; i < m; i++ {
		b[i] = a[i]
	}
	return b
}

// float64ToComplex128 converts a float64 array to a complex128 array
func Float64_to_complex128(a []float64) []complex128 {
	n := len(a)
	b := make([]complex128, n)
	for i, v := range a {
		b[i] = complex(v, 0)
	}
	return b
}

func Real_part_of_complex_array(complexArr []complex128) []float64 {
	floatArr := make([]float64, len(complexArr))
	for i, c := range complexArr {
		floatArr[i] = real(c)
	}
	return floatArr
}