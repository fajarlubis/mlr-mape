package model

import "time"

// struct untuk menampung data file csv
type Result struct {
	Items []*ResultItem
	Type  string
	MAPE  float64
}

// struct untuk menampung data item-item per bulan file csv
type ResultItem struct {
	Year       int
	Month      int
	Prediction float64
}

// struct untuk mapping data tahunan eksport ke xlsx file
type YearlyData struct {
	Year int
	Data map[string]map[time.Month]float64
}
