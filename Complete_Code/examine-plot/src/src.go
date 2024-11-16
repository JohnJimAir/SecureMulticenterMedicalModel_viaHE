package src

import (
	"local.com/plain-math/statistics"
)


func Scores_which(yPred, yTrue_0, yTrue_1, labels_allowed []int, label_positive int) ([]float64) {
	
	var yTrue []int
	if statistics.Accuracy(yPred, yTrue_0, labels_allowed) > statistics.Accuracy(yPred, yTrue_1, labels_allowed) {
		yTrue = yTrue_0
	} else {
		yTrue = yTrue_1
	}

	return statistics.Scores(yPred, yTrue, labels_allowed, label_positive)
}

func AUC_which(probs []float64, yPred, yTrue_0, yTrue_1, labels_allowed []int, label_positive int) (auc float64, tpr_record, fpr_record []float64) {

	var yTrue []int
	if statistics.Accuracy(yPred, yTrue_0, labels_allowed) > statistics.Accuracy(yPred, yTrue_1, labels_allowed) {
		yTrue = yTrue_0
	} else {
		yTrue = yTrue_1
	}

	return statistics.AUC(probs, yTrue, labels_allowed, label_positive)
}
