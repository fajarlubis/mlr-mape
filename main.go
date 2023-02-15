package main

import (
	"encoding/csv"
	"os"
	"strconv"
)

// func main() {
// 	// Load the dataset from a CSV file
// 	file, err := os.Open("energy_loss.csv")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer file.Close()

// 	reader := csv.NewReader(file)
// 	records, err := reader.ReadAll()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Extract the input and output data from the dataset
// 	var x [][]float64
// 	var y []float64

// 	for _, record := range records {
// 		var input []float64
// 		for i := 0; i < len(record); i++ {
// 			val, err := strconv.ParseFloat(record[i], 64)
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 			input = append(input, val)
// 		}
// 		x = append(x, input)

// 		output, err := strconv.ParseFloat(record[len(record)-1], 64)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		y = append(y, output)
// 	}

// 	// Train the linear regression model
// 	r := new(regression.Regression)
// 	r.SetObserved("Energy loss every year")
// 	r.SetVar(0, "month")
// 	r.SetVar(1, "year")
// 	r.SetVar(2, "energy_loss_1")
// 	r.SetVar(3, "energy_loss_2")

// 	for i, xi := range x {
// 		log.Println(y[i], xi)
// 		r.Train(regression.DataPoint(y[i], xi))
// 	}

// 	r.Run()

// 	// Use the trained model to make predictions for the next 5 years
// 	var predictions []float64
// 	for year := 2019; year <= 2023; year++ {
// 		for month := 1; month <= 12; month++ {
// 			input := []float64{float64(month), float64(year)}
// 			prediction, err := r.Predict(input)
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 			predictions = append(predictions, prediction)
// 		}
// 	}

// 	// Calculate the mean absolute percentage error of the predictions
// 	var actual []float64
// 	var absoluteErrors []float64
// 	actual = append(actual, y...)

// 	for i, value := range predictions {
// 		absoluteError := 100 * (actual[i] - value) / actual[i]
// 		absoluteErrors = append(absoluteErrors, absoluteError)
// 	}

// 	sum := 0.0
// 	for _, error := range absoluteErrors {
// 		sum += error
// 	}
// 	meanAbsolutePercentageError := sum / float64(len(absoluteErrors))

// 	// Print the predictions and the mean absolute percentage error
// 	fmt.Println("Predictions:")
// 	for i, prediction := range predictions {
// 		fmt.Printf("%d-%02d: %.2f\n", (i/12)+2019, (i%12)+1, prediction)
// 	}

// 	fmt.Printf("Mean absolute percentage error: %.2f%%\n", meanAbsolutePercentageError)
// }

// =========================================================================================

// func main() {
// 	// Load the dataset from a CSV file
// 	data, err := readCSV("energy_loss2.csv")
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Extract the input and output data from the dataset
// 	var x [][]float64
// 	var y []float64

// 	for _, record := range records {
// 		var input []float64
// 		for i := 0; i < len(record)-1; i++ {
// 			val, err := strconv.ParseFloat(record[i], 64)
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 			input = append(input, val)
// 		}
// 		x = append(x, input)

// 		output, err := strconv.ParseFloat(record[len(record)-1], 64)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		y = append(y, output)
// 	}

// 	// Train the linear regression model
// 	r := new(regression.Regression)
// 	r.SetObserved("energy_loss")
// 	r.SetVar(0, "month")
// 	r.SetVar(1, "year")

// 	for i, xi := range x {
// 		r.Train(regression.DataPoint(y[i], xi))
// 	}

// 	r.Run()

// 	// Use the trained model to make predictions for the next 5 years
// 	var predictions []float64
// 	for year := 2019; year <= 2023; year++ {
// 		for month := 1; month <= 12; month++ {
// 			input := []float64{float64(month), float64(year)}
// 			prediction, err := r.Predict(input)
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 			predictions = append(predictions, prediction)
// 		}
// 	}

// 	// Calculate the mean absolute percentage error of the predictions
// 	var actual []float64
// 	var absoluteErrors []float64
// 	actual = append(actual, y...)

// 	for i, value := range predictions {
// 		absoluteError := 100 * (actual[i] - value) / actual[i]
// 		absoluteErrors = append(absoluteErrors, absoluteError)
// 	}

// 	sum := 0.0
// 	for _, error := range absoluteErrors {
// 		sum += error
// 	}
// 	meanAbsolutePercentageError := sum / float64(len(absoluteErrors))

// 	// Print the predictions and the mean absolute percentage error
// 	fmt.Println("Predictions:")
// 	for i, prediction := range predictions {
// 		fmt.Printf("%d-%02d: %.2f\n", (i/12)+2019, (i%12)+1, prediction)
// 	}

// 	fmt.Printf("Mean absolute percentage error: %.2f%%\n", meanAbsolutePercentageError)
// }

// =========================================================================================

type EnergyLoss struct {
	Month       int
	Year        int
	EnergyLoss1 float64
	EnergyLoss2 float64
	EnergyLoss3 float64
}

// func main() {
// 	// Read the CSV file
// 	data, err := readCSV("energy_loss2.csv")
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Create the regression model
// 	r := new(regression.Regression)
// 	r.SetObserved("EnergyLoss1")
// 	r.SetVar(0, "Month")
// 	r.SetVar(1, "Year")
// 	r.SetVar(2, "EnergyLoss1")
// 	r.SetVar(3, "EnergyLoss2")
// 	r.SetVar(4, "EnergyLoss3")

// 	// Add the data to the regression model
// 	for _, row := range data {
// 		r.Train(
// 			regression.DataPoint(20,
// 				[]float64{
// 					float64(row.Month),
// 					float64(row.Year),
// 					row.EnergyLoss1,
// 					row.EnergyLoss2,
// 					row.EnergyLoss3,
// 				},
// 			),
// 		)
// 	}

// 	// Train the regression model
// 	r.Run()

// 	// Predict future values
// 	var predictions []EnergyLoss
// 	for year := 2023; year <= 2027; year++ {
// 		for month := 1; month <= 12; month++ {
// 			x := []float64{float64(month), float64(year)}
// 			x = append(x, r.GetCoeffs()...)
// 			y, err := r.Predict(x)
// 			if err != nil {
// 				panic(err)
// 			}
// 			predictions = append(predictions, EnergyLoss{
// 				Month:       month,
// 				Year:        year,
// 				EnergyLoss1: y,
// 				EnergyLoss2: y,
// 				EnergyLoss3: y,
// 			})
// 		}
// 	}

// 	// Calculate Mean Absolute Percentage Error
// 	var mape float64
// 	for i, row := range data {
// 		predictedRow := predictions[i]
// 		mape += math.Abs((predictedRow.EnergyLoss1-row.EnergyLoss1)/row.EnergyLoss1) +
// 			math.Abs((predictedRow.EnergyLoss2-row.EnergyLoss2)/row.EnergyLoss2) +
// 			math.Abs((predictedRow.EnergyLoss3-row.EnergyLoss3)/row.EnergyLoss3)
// 	}
// 	mape /= float64(len(data) * 3)

// 	// Print the predictions and the mean absolute percentage error
// 	fmt.Println("Predictions:")
// 	for i, prediction := range predictions {
// 		fmt.Printf("%d-%02d: %.2f\n", (i/12)+2019, (i%12)+1, prediction.EnergyLoss3)
// 	}

// 	fmt.Printf("Mean Absolute Percentage Error: %.2f%%\n", mape*100)
// }

// func main() {
// 	r := new(regression.Regression)
// 	r.SetObserved("Murders per annum per 1,000,000 inhabitants")
// 	r.SetVar(0, "Inhabitants")
// 	r.SetVar(1, "Percent with incomes below $5000")
// 	r.SetVar(2, "Percent unemployed")
// 	r.Train(
// 		regression.DataPoint(11.2, []float64{587000, 16.5, 6.2}),
// 		regression.DataPoint(13.4, []float64{643000, 20.5, 6.4}),
// 		regression.DataPoint(40.7, []float64{635000, 26.3, 9.3}),
// 		regression.DataPoint(5.3, []float64{692000, 16.5, 5.3}),
// 		regression.DataPoint(24.8, []float64{1248000, 19.2, 7.3}),
// 		regression.DataPoint(12.7, []float64{643000, 16.5, 5.9}),
// 		regression.DataPoint(20.9, []float64{1964000, 20.2, 6.4}),
// 		regression.DataPoint(35.7, []float64{1531000, 21.3, 7.6}),
// 		regression.DataPoint(8.7, []float64{713000, 17.2, 4.9}),
// 		regression.DataPoint(9.6, []float64{749000, 14.3, 6.4}),
// 		regression.DataPoint(14.5, []float64{7895000, 18.1, 6}),
// 		regression.DataPoint(26.9, []float64{762000, 23.1, 7.4}),
// 		regression.DataPoint(15.7, []float64{2793000, 19.1, 5.8}),
// 		regression.DataPoint(36.2, []float64{741000, 24.7, 8.6}),
// 		regression.DataPoint(18.1, []float64{625000, 18.6, 6.5}),
// 		regression.DataPoint(28.9, []float64{854000, 24.9, 8.3}),
// 		regression.DataPoint(14.9, []float64{716000, 17.9, 6.7}),
// 		regression.DataPoint(25.8, []float64{921000, 22.4, 8.6}),
// 		regression.DataPoint(21.7, []float64{595000, 20.2, 8.4}),
// 		regression.DataPoint(25.7, []float64{3353000, 16.9, 6.7}),
// 	)
// 	r.Run()

// 	fmt.Printf("Regression formula:\n%v\n", r.Formula)
// 	fmt.Printf("Regression:\n%s\n", r)
// }

func readCSV(filename string) ([]EnergyLoss, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 5
	reader.TrimLeadingSpace = true

	lines, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var data []EnergyLoss
	for i, line := range lines {
		if i == 0 {
			continue
		}
		month, _ := strconv.Atoi(line[0])
		year, _ := strconv.Atoi(line[1])
		energyLoss1, _ := strconv.ParseFloat(line[2], 64)
		energyLoss2, _ := strconv.ParseFloat(line[3], 64)
		energyLoss3, _ := strconv.ParseFloat(line[4], 64)

		data = append(data, EnergyLoss{
			Month:       month,
			Year:        year,
			EnergyLoss1: energyLoss1,
			EnergyLoss2: energyLoss2,
			EnergyLoss3: energyLoss3,
		})
	}

	return data, nil
}
