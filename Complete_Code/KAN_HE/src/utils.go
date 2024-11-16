package src

func Interval(endpoint []float64) (interval [][]float64) {

	length := len(endpoint)
	interval = make([][]float64, length)
	for i:=0;i<length;i++ {
		interval[i] = []float64{-endpoint[i], endpoint[i]}
	}
	return interval
}