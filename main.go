package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/sajari/regression"
)

type Result struct {
	Items []*ResultItem
	MAPE  float64
}

type ResultItem struct {
	Year       int
	Month      int
	Prediction float64
}

func main() {
	// Load the dataset from a CSV file
	file, err := os.Open("energy_loss3.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	var (
		res   = make([]*Result, 0)
		items = make([]*ResultItem, 0)
	)

	for label := 0; label <= 2; label++ { // there is 11 data label
		// Extract the input and output data from the dataset
		var (
			x [][]float64
			y []float64
		)

		for _, record := range records {
			var input []float64
			for i := 0; i < len(record)-3; i++ {
				val, err := strconv.ParseFloat(record[i], 64)
				if err != nil {
					log.Fatal(err)
				}
				input = append(input, val)
			}
			x = append(x, input)

			yOut, err := strconv.ParseFloat(record[label+2], 64)
			if err != nil {
				log.Fatal(err)
			}
			y = append(y, yOut)
		}

		// Train the linear regression model
		r := new(regression.Regression)
		r.SetObserved("Energy loss every year")
		r.SetVar(0, "month")
		r.SetVar(1, "year")

		for i, xi := range x {
			r.Train(
				regression.DataPoint(y[i], xi),
			)
		}

		r.Run()

		// fmt.Printf("Regression formula:\n%v\n", r.Formula)
		// fmt.Printf("Regression:\n%s\n", r)

		// Use the trained model to make predictions for the next 5 years
		var predictions []float64
		for year := 2019; year <= 2023; year++ {
			for month := 1; month <= 12; month++ {
				input := []float64{float64(month), float64(year)}
				prediction, err := r.Predict(input)
				if err != nil {
					log.Fatal(err)
				}
				predictions = append(predictions, prediction)
			}
		}

		// Calculate the mean absolute percentage error of the predictions
		var actual []float64
		var absoluteErrors []float64
		actual = append(actual, y...)

		for i, value := range predictions {
			absoluteError := 100 * (actual[i] - value) / actual[i]
			absoluteErrors = append(absoluteErrors, absoluteError)
		}

		sum := 0.0
		for _, error := range absoluteErrors {
			sum += error
		}
		meanAbsolutePercentageError := sum / float64(len(absoluteErrors))

		// Print the predictions and the mean absolute percentage error
		// fmt.Println("Predictions:")
		for i, prediction := range predictions {
			items = append(items, &ResultItem{
				Year:       (i / 12) + 2019,
				Month:      (i % 12) + 1,
				Prediction: prediction,
			})

			// fmt.Printf("%d-%02d: %.2f\n", (i/12)+2019, (i%12)+1, prediction)
		}

		// fmt.Printf("Mean absolute percentage error: %.2f%%\n", meanAbsolutePercentageError)

		res = append(res, &Result{
			Items: items,
			MAPE:  meanAbsolutePercentageError,
		})
	}

	for _, v := range res {
		for _, w := range v.Items {
			fmt.Printf("%d-%02d: %.2f\n", w.Year, w.Month, w.Prediction)
		}
		fmt.Printf("Mean absolute percentage error: %.2f%%\n", v.MAPE)
	}
}
