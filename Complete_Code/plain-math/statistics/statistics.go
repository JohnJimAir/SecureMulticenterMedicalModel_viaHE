package statistics

import "local.com/plain-math/vecop"


func Scores(yPred, yTrue, labels_allowed []int, label_positive int) ([]float64) {

	accuracy := Accuracy(yPred, yTrue, labels_allowed)
	precision := Precision(yPred, yTrue, labels_allowed, label_positive)
	recall := Recall(yPred, yTrue, labels_allowed, label_positive)
	f1score := (precision + recall) / 2

	return []float64{accuracy, precision, recall, f1score}
}

func Accuracy(yPred, yTrue, labels_allowed []int) float64 {
	
	Validate(yPred, yTrue, labels_allowed)

	return vecop.CommonRatio_int(yPred, yTrue)
}

func Precision(yPred, yTrue, labels_allowed []int, label_positive int) float64 {

	tp := TP(yPred, yTrue, labels_allowed, label_positive)
	fp := FP(yPred, yTrue, labels_allowed, label_positive)
	return float64(tp) / float64(tp + fp)
}

func Recall(yPred, yTrue, labels_allowed []int, label_positive int) float64 {

	tp := TP(yPred, yTrue, labels_allowed, label_positive)
	fn := FN(yPred, yTrue, labels_allowed, label_positive)
	return float64(tp) / float64(tp + fn)
}


func TP(yPred, yTrue, labels_allowed []int, label_positive int) int {

	Validate(yPred, yTrue, labels_allowed)

	count := 0
	for i:=0; i<len(yPred); i++ {
		if yPred[i] == label_positive && yTrue[i] == label_positive {
			count++
		}
	}
	return count
}

func FP(yPred, yTrue, labels_allowed []int, label_positive int) int {

	Validate(yPred, yTrue, labels_allowed)

	count := 0
	for i:=0; i<len(yPred); i++ {
		if yPred[i] == label_positive && yTrue[i] != label_positive {
			count++
		}
	}
	return count
}

func TN(yPred, yTrue, labels_allowed []int, label_positive int) int {

	Validate(yPred, yTrue, labels_allowed)

	count := 0
	for i:=0; i<len(yPred); i++ {
		if yPred[i] != label_positive && yTrue[i] != label_positive {
			count++
		}
	}
	return count
}

func FN(yPred, yTrue, labels_allowed []int, label_positive int) int {

	Validate(yPred, yTrue, labels_allowed)

	count := 0
	for i:=0; i<len(yPred); i++ {
		if yPred[i] != label_positive && yTrue[i] == label_positive {
			count++
		}
	}
	return count
}

func Validate(yPred, yTrue, labels_allowed []int) () {

	if len(yPred) != len(yTrue) || len(yPred) == 0 || len(yTrue) == 0 {
		panic("length wrong.")
	}
	if !vecop.CheckAllowedElements(yPred, labels_allowed) || !vecop.CheckAllowedElements(yTrue, labels_allowed) {
		panic("not 0 or 1.")
	}
}