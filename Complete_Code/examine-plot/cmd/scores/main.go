package main

import (
	"local.com/examine-plot/src"
	"local.com/numerical-read-print/nprint"
	"local.com/numerical-read-print/nread"
	"local.com/plain-math/matop"
	"local.com/plain-math/statistics"
	"local.com/plain-math/vecop"
)

func main() {
	
	labels_allowed := []int{0, 1}
	label_positive := 1
	disease := "sepsis"
	domain := "cipher"
	model := "KAN"
	filename_result := "../../data/result/" + disease + "_" + model + "_" + domain + ".csv"

	var yTrue [][]int
	var result [][]float64
	if model == "KAN" {   
		filename_label := "../../data/label/label_" + disease + "_KAN.txt"
		yTrue = nread.Int_from_float64(matop.Transpose(nread.M_Blank(filename_label)))
		inverse := make([]int, len(yTrue[0]))
		for i:=0; i<len(inverse); i++ {
			inverse[i] = 1 - yTrue[0][i]
		}
		yTrue = append(yTrue, inverse)

		if disease == "sepsis" {
			result = nread.M_Comma(filename_result)
			for i:=0; i<len(result); i++ {
				result[i] = append(result[i], 0.5)
			}
		} else {
			result = nread.M_Comma(filename_result)
		}
	} else {
		filename_label := "../../data/label/label_" + disease + ".txt"
		yTrue = nread.Int_from_float64(matop.Transpose(nread.M_Blank(filename_label)))
		result = nread.M_Comma(filename_result)
	}

	result_T := matop.Transpose(result)
	yPred := vecop.Bigger(result_T[0], result_T[1])

	scores := src.Scores_which(yPred, yTrue[0], yTrue[1], labels_allowed, label_positive)

	res_softmax := matop.Softmax_mat(result)
	res_softmax_T := matop.Transpose(res_softmax)
	auc, tpr_record, fpr_record := src.AUC_which(res_softmax_T[0], yPred, yTrue[0], yTrue[1], labels_allowed, label_positive)
	scores = append(scores, auc)
	nprint.Print_Vector(scores, 6, len(scores))
	nprint.Print_Vector(tpr_record, 6, len(tpr_record))
	nprint.Print_Vector(fpr_record, 6, len(fpr_record))
	if err := statistics.PlotROC(tpr_record, fpr_record, "./ROC.png"); err != nil {
		panic("plot error.")
	}

}