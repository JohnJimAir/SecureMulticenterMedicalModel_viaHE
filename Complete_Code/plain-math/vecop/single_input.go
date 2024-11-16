package vecop

import "math"

func Sum(vec []float64) (result float64) {

	if len(vec) == 0 {
		panic("length wrong.")
	}

	result = 0
	for i:=0; i<len(vec); i++ {
		result += vec[i]
	}
	return result
}

func Product(vec []float64) (result float64) {

	if len(vec) == 0 {
		panic("length wrong.")
	}

	result = 1.0
	for i:=0; i<len(vec); i++ {
		result *= vec[i]
	}
	return result
}

func Contains(vec []int, element int) bool {
	
	if len(vec) == 0 {
		panic("length wrong.")
	}

	for i:=0; i<len(vec); i++ {
		if vec[i] == element {
			return true
		}
	}
	return false
}

func CountSpecificElements(vec, elements []int) (count int) {

	if len(vec) == 0 || len(elements) == 0 {
		panic("length wrong.")
	}

	count = 0
	for i:=0; i<len(vec); i++ {
		if Contains(elements, vec[i]) {
			count++
		}
	}
	return count
}

func CheckAllowedElements(vec, elements []int) bool {

	if len(vec) == 0 || len(elements) == 0 {
		panic("length wrong.")
	}

	for i:=0; i<len(vec); i++ {
		if Contains(elements, vec[i]) == false {
			return false
		}
	}

	return true
}

func CheckAllowedRange(vec []float64, interval []float64, epsilon float64) bool {

	if len(vec) == 0 || len(interval) != 2 {
		panic("length wrong.")
	}

	for _, v := range vec {
		if v < (interval[0] - epsilon) || v > (interval[1] + epsilon) {
			return false
		}
	}
	return true
}

func Repeat(value float64, count int) (vec []float64) {

	vec = make([]float64, count)
	for i:=0; i<count; i++ {
		vec[i] = value
	}
	return vec
}

func Softmax(vec []float64) []float64 {

	if len(vec) <= 1 {
		panic("length wrong.")
	}
	
	expSum := 0.0
	for _, v := range vec {
		expSum += math.Exp(v)
	}

	result := make([]float64, len(vec))
	for i, v := range vec {
		result[i] = math.Exp(v) / expSum
	}
	return result
}



func Max() () {}

func Index_Max() () {}