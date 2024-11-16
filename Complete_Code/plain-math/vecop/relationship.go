package vecop

import "math"

func Bigger(vec0, vec1 []float64) (result []int) {

	if len(vec0) != len(vec1) || len(vec0) == 0 || len(vec1) == 0 {
		panic("length wrong.")
	}

	result = make([]int, len(vec0))
	for i:=0; i<len(vec0); i++ {
		if vec0[i] > vec1[i] {
			result[i] = 1
		} else {
			result[i] = 0
		}
	}
	return result
}

func Smaller(vec0, vec1 []float64) (result []int) {
	
	result = Bigger(vec0, vec1)

	for i:=0; i<len(result); i++ {
		result[i] = 1 - result[i]
	}
	return result
}

func CommonRatio_float(vec0, vec1 []float64, epsilon float64) float64 {

	if len(vec0) != len(vec1) || len(vec0) == 0 || len(vec1) == 0 {
		panic("length wrong.")
	}
	if epsilon <= 0 {
		panic("epsilon must be bigger than 0.")
	}

	commonCount := 0
    length := len(vec0)

    for i := 0; i < length; i++ {
        if math.Abs(vec0[i] - vec1[i]) <= epsilon {
            commonCount++
        }
    }
    return float64(commonCount) / float64(length)
}

func CommonRatio_int(vec0, vec1 []int) float64 {
    
	if len(vec0) != len(vec1) || len(vec0) == 0 || len(vec1) == 0 {
		panic("length wrong.")
	}

	commonCount := 0
    length := len(vec0)

    for i := 0; i < length; i++ {
        if vec0[i] == vec1[i] {
            commonCount++
        }
    }
    return float64(commonCount) / float64(length)
}

