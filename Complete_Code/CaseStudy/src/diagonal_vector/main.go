package main

import (
	"errors"
	"fmt"
)

// Function to compute diagonal vectors
func diagonalVectors(A [][]float64) ([][]float64, error) {
    m := len(A)
    n := len(A[0])
    
    // Check if the matrix is square
    if m != n {
        return nil, errors.New("the input must be a square matrix")
    }

    // Initialize the result matrix
    d := make([][]float64, n)
    for i := range d {
        d[i] = make([]float64, n)
    }

    // Compute the diagonal vectors
    for i := 0; i < n; i++ {
        for j := 0; j < n; j++ {
            if i+j < n {
                d[i][j] = A[j][i+j]
            } else {
                d[i][j] = A[j][(i+j)%n]
            }
        }
    }

    return d, nil
}

func main() {
    // Example usage
    A := [][]float64{
        {1, 2, 3},
        {4, 5, 6},
        {7, 8, 9},
    }

    d, err := diagonalVectors(A)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    for _, row := range d {
        fmt.Println(row)
    }
}
