package statistics

import (
	"sort"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"local.com/plain-math/vecop"
)

type Prediction struct {
	prob float64
	label int
}

func AUC(probs []float64, labels, label_allowed []int, label_positive int) (auc float64, tpr_record, fpr_record []float64) {

	if len(probs) != len(labels) || len(probs) == 0 {
		panic("length wrong.")
	}
	if !vecop.CheckAllowedRange(probs, []float64{0.0, 1.0}, 0.0001) {
		panic("range of probs wrong.")
	}

	count_pos := vecop.CountSpecificElements(labels, []int{label_positive})
	if count_pos == len(labels) || count_pos == 0 {
		panic("labels wrong.")
	}
	length := len(probs)
	count_neg := length - count_pos

	predictions := make([]Prediction, length)
	for i:=0; i<length; i++ {
		predictions[i] = Prediction{probs[i], labels[i]}
	}
	sort.Slice(predictions, func(i, j int) bool {
		return predictions[i].prob > predictions[j].prob
	})

	tpr_record = append(tpr_record, 0.0)
	fpr_record = append(fpr_record, 0.0)
	var tp, fp, tpr, fpr, fpr_pre float64
	for _, prediction := range predictions {

		if prediction.label == label_positive {
			tp++
		} else {
			fp++
		}

		tpr = float64(tp) / float64(count_pos)
		fpr = float64(fp) / float64(count_neg)
		tpr_record = append(tpr_record, tpr)
		fpr_record = append(fpr_record, fpr)

		auc += (fpr - fpr_pre) * tpr

		fpr_pre = fpr
	}

	return auc, tpr_record, fpr_record
}

func PlotROC(tpr []float64, fpr []float64, filename string) error {

	p := plot.New()

	p.Title.Text = "ROC Curve"
	p.X.Label.Text = "False Positive Rate"
	p.Y.Label.Text = "True Positive Rate"

	pts := make(plotter.XYs, len(tpr))
	for i := range tpr {
		pts[i].X = fpr[i]
		pts[i].Y = tpr[i]
	}

	line, err := plotter.NewLine(pts)
	if err != nil {
		return err
	}
	p.Add(line)

	if err := p.Save(6*vg.Inch, 3*vg.Inch, filename); err != nil {
		return err
	}
	return nil
}