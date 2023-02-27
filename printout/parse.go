package printout

import "mlr-mape/model"

type Result struct {
	Items map[int]map[int]float64 // nested map to store prediction data
	Type  string
	MAPE  float64
}

func Parse(items []*model.ResultItem, resultType string, mape float64) *Result {
	result := &Result{
		Items: make(map[int]map[int]float64),
		Type:  resultType,
		MAPE:  mape,
	}

	// Populate nested map with prediction data from ResultItems
	for _, item := range items {
		if _, ok := result.Items[item.Year]; !ok {
			result.Items[item.Year] = make(map[int]float64)
		}
		result.Items[item.Year][item.Month] = item.Prediction
	}

	return result
}
