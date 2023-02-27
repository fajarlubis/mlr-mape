package main

import (
	"encoding/csv"
	"log"
	"mlr-mape/model"
	"mlr-mape/printout"
	"os"
	"strconv"

	"github.com/sajari/regression"
)

func main() {
	// memuat dataset dari file csv
	file, err := os.Open("training-data.csv")
	if err != nil {
		log.Fatal(err)
	}
	// fungsi defer file.Close() untuk menunda close buffer dari os.Open setelah selesai mengeksekusi function keseluruhan
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll() // proses membaca isi file csv
	if err != nil {
		log.Fatal(err)
	}

	var (
		// definisi variable baru dari struct diatas
		res   = make([]*model.Result, 0)
		items = make([]*model.ResultItem, 0)
	)

	for label := 0; label <= 10; label++ { // ada 11 data label mulai dari kWh penerimaan sampai susut non-teknis, di buat 10 karena index dimulai dari 0

		// ekstrak input data dari file csv
		var (
			// definisi variable x, y dan kWhType
			x       [][]float64
			y       []float64
			kWhType string
		)

		for _, record := range records { // loop setiap data yang sudah di tampung di variable yang sudah di definisikan sebelumnya
			var input []float64                   // definisi variable baru untuk nilai input float64
			for i := 0; i < len(record)-11; i++ { // loop 2 kolom pertama yang berisikan bulan dan tahun
				val, err := strconv.ParseFloat(record[i], 64) // konversi isi data ke bentuk float64
				if err != nil {
					log.Fatal(err) // jika bulan dan tahun bukan dalam bentuk angka maka return error
				}
				input = append(input, val) // jika bulan dan tahun dalam bentuk angka maka sematkan ke dalam variable input
			}
			x = append(x, input) // setelah loop 1 row selesai maka variable x diisi dengan hasil dari loop bulan dan tahun

			yOut, err := strconv.ParseFloat(record[label+2], 64) // yOut adalah isi dari setiap cell yang ada di dalam file csv, dimulai dari label+2 yang berarti jika index sekarang 0 maka dimulai dari kolom 2 (2+0)
			if err != nil {
				log.Fatal(err) // jika konversi nilai gagal maka return error
			}
			y = append(y, yOut) // setelah loop 1 row selesai maka variable y diisi dengan nilai dari setiap cell yang ada di dalam file csv

			// switch label + 2 menentukan jenis/label dari setiap nilai, label + 2 menandakan posisi kolom saat ini, +2 karena 2 kolom pertama bukan berisi nilai melainkan bulan dan tahun
			// jika nilai sesuai dengan masing-masing case maka variable kWhType akan berisi dengan nama dari masing-masing label di setiap nilainya
			switch label + 2 {
			case 2:
				kWhType = "kWh Penerimaan"
			case 3:
				kWhType = "kWh Penjualan"
			case 4:
				kWhType = "Pemakaian Sendiri"
			case 5:
				kWhType = "Kirim ke Unit Lain"
			case 6:
				kWhType = "Susut Teknis JTM"
			case 7:
				kWhType = "Susut Teknis JTR"
			case 8:
				kWhType = "Susut Teknis Trafo"
			case 9:
				kWhType = "Susut Teknis SR"
			case 10:
				kWhType = "Susut Total"
			case 11:
				kWhType = "Susut Teknis"
			case 12:
				kWhType = "Susut Non-Teknis"
			}
		}

		// training linear regression model
		r := new(regression.Regression) // inisialisasi regressi baru
		r.SetObserved("Nilai kWh")      // set label jenis variable yang akan di observasi/diproses
		r.SetVar(0, "month")            // set variable month(bulan) dengan index pertama 0
		r.SetVar(1, "year")             // set variable year(tahun) dengan index kedua 1

		for i, xi := range x { // loop berdasarkan sumbu x dengan data-data didalamnya
			r.Train( // mempersiapkan training data sesuai nilai x dan y yang sudah diberikan
				regression.DataPoint(y[i], xi), // menambahkan datapoint ke dalam proses training
			)
		}

		r.Run() // start training data sesuai dengan instruksi dan nilai yang telah diberikan

		// fmt.Printf("Regression formula:\n%v\n", r.Formula)
		// fmt.Printf("Regression:\n%s\n", r)

		// fmt.Printf("R^2:\n%v\n", r.R2) // print hasil R^2 dari hasil training ke layar

		// inisialisasi variable untuk menampung data prediksi untuk 5 tahun kedepan
		var predictions []float64
		for year := 2021; year <= 2025; year++ { // loop jarak tahun yang akan di prediksi 2021 - 2025
			for month := 1; month <= 12; month++ { // loop setiap bulan 1 - 12
				input := []float64{float64(month), float64(year)} // mengambil data bulan dan tahun dan mengisi kedalam variable slice baru
				prediction, err := r.Predict(input)               // prediksi sesuai data yang sudah di training
				if err != nil {
					log.Fatal(err) // jika error maka return dan keluarkan error di layar
				}
				predictions = append(predictions, prediction) // menyematkan hasil prediksi ke dalam variable predictions
			}
		}

		// inisialisasi menghitung Mean Absolute Percentage Error (MAPE)
		var actual []float64
		var absoluteErrors []float64
		actual = append(actual, y...) // mengambil data aktual hasil training dan menyematkan ke dalam variable baru actual

		for i, value := range predictions { // loop data prediksi
			absoluteError := 100 * (actual[i] - value) / actual[i] // mulai menghitung MAPE sesuai rumus
			absoluteErrors = append(absoluteErrors, absoluteError) // menyematkan data hasil MAPE ke dalam variable absoluteErrors
		}

		sum := 0.0                            // nilai sum awal 0.0 default
		for _, errx := range absoluteErrors { // loop absoluteErrors dan menambahkan total absolute error
			sum += errx
		}
		meanAbsolutePercentageError := sum / float64(len(absoluteErrors))

		// menyematkan data hasil prediksi dan MAPE ke dalam variable predictions untuk dapat di proses ke tahap selanjutnya
		for i, prediction := range predictions {
			items = append(items, &model.ResultItem{
				Year:       (i / 12) + 2021, // tahun mulai yang akan di prediksi
				Month:      (i % 12) + 1,    // bulan dalam angka
				Prediction: prediction,
			})

			// fmt.Printf("%d-%02d: %.2f\n", (i/12)+2019, (i%12)+1, prediction)
		}

		// fmt.Printf("Mean absolute percentage error: %.2f%%\n", meanAbsolutePercentageError)

		// menyematkan data setiap item ke dalam variable result
		res = append(res, &model.Result{
			Items: items,
			Type:  kWhType,
			MAPE:  meanAbsolutePercentageError,
		})
	}

	// ekspor hasil training ke dalam bentuk file xlsx
	if err := printout.ExportLegacy(res); err != nil {
		log.Fatalln(err) // jika error maka data tidak akan muncul di file result.xlsx
	}

	// predictions := make([]YearlyData, 0)

	// // var year int
	// for _, v := range res {
	// 	// var currentYear int
	// 	// fmt.Println("===================================================================>", v.MAPE, v.Type, "<===================================================================")

	// 	// for _, w := range v.Items {
	// 	// 	currentYear = w.Year

	// 	// 	if currentYear != year {
	// 	// 		year = currentYear
	// 	// 		fmt.Println("=================================================================== --- ===================================================================")
	// 	// 	}

	// 	// 	fmt.Println(time.Month(w.Month), w.Prediction, "Type", v.Type, "TAHUN", w.Year)
	// 	// }

	// 	// ----------------------------------

	// 	switch v.Type {
	// 	case "kWh Penerimaan":
	// 		monthData := make(map[time.Month]float64, 0)
	// 		yearData := make(map[string]map[time.Month]float64, 0)

	// 		var currentYear int
	// 		for _, w := range v.Items {
	// 			if currentYear != w.Year {
	// 				predictions = append(predictions, YearlyData{
	// 					Year: w.Year,
	// 					Data: yearData,
	// 				})

	// 				currentYear = w.Year
	// 				monthData = make(map[time.Month]float64)
	// 				yearData = make(map[string]map[time.Month]float64, 0)

	// 				fmt.Println("=================================================================== --- ===================================================================")
	// 				// continue
	// 			}

	// 			monthData[time.Month(w.Month)] = w.Prediction
	// 			yearData[v.Type] = monthData

	// 			fmt.Println(time.Month(w.Month), w.Prediction, "Type", v.Type, "TAHUN", w.Year)
	// 		}

	// 		fmt.Print("\n\n")
	// 	case "kWh Penjualan":
	// 		monthData := make(map[time.Month]float64, 0)
	// 		yearData := make(map[string]map[time.Month]float64, 0)

	// 		var currentYear int
	// 		for _, w := range v.Items {
	// 			if currentYear != w.Year {
	// 				predictions = append(predictions, YearlyData{
	// 					Year: w.Year,
	// 					Data: yearData,
	// 				})

	// 				currentYear = w.Year
	// 				monthData = make(map[time.Month]float64)
	// 				yearData = make(map[string]map[time.Month]float64, 0)

	// 				fmt.Println("=================================================================== --- ===================================================================")
	// 				// continue
	// 			}

	// 			monthData[time.Month(w.Month)] = w.Prediction
	// 			yearData[v.Type] = monthData

	// 			fmt.Println(time.Month(w.Month), w.Prediction, "Type", v.Type, " TAHUN", w.Year)
	// 		}

	// 		fmt.Print("\n\n")
	// 	}
	// }

	// for _, v := range predictions {
	// 	for j, w := range v.Data {
	// 		fmt.Print(v.Year, j, w, "\n\n")
	// 	}
	// }

	// // for _, v := range res {
	// // 	var currentYear int
	// // 	for _, w := range v.Items {
	// // 		currentYear = w.Year

	// // 		monthData[time.Month(w.Month)] = w.Prediction
	// // 		yearData[v.Type] = monthData
	// // 	}

	// // 	predictions = append(predictions, YearlyData{
	// // 		Year: currentYear,
	// // 		Data: yearData,
	// // 	})
	// // }

	// if err := writeResult(nil, predictions); err != nil {
	// 	log.Println(err)
	// }
}
