package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sajari/regression"
	"github.com/xuri/excelize/v2"
)

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
		res   = make([]*Result, 0)
		items = make([]*ResultItem, 0)
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

		fmt.Printf("R^2:\n%v\n", r.R2) // print hasil R^2 dari hasil training ke layar

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
			items = append(items, &ResultItem{
				Year:       (i / 12) + 2021, // tahun mulai yang akan di prediksi
				Month:      (i % 12) + 1,    // bulan dalam angka
				Prediction: prediction,
			})

			// fmt.Printf("%d-%02d: %.2f\n", (i/12)+2019, (i%12)+1, prediction)
		}

		// fmt.Printf("Mean absolute percentage error: %.2f%%\n", meanAbsolutePercentageError)

		// menyematkan data setiap item ke dalam variable result
		res = append(res, &Result{
			Items: items,
			Type:  kWhType,
			MAPE:  meanAbsolutePercentageError,
		})
	}

	// ekspor hasil training ke dalam bentuk file xlsx
	if err := export2(res); err != nil {
		log.Fatalln(err) // jika error maka data tidak akan muncul di file result.xlsx
	}
}

func export2(data []*Result) error {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	var (
		defaultSheet = "Sheet1"
	)

	// create a new sheet
	index, err := f.NewSheet(defaultSheet)
	if err != nil {
		return err
	}

	var kWhType string
	rowStartIdx := 1
	for _, res := range data {
		kWhType = res.Type

		for _, v := range res.Items {
			f.SetCellValue(defaultSheet, fmt.Sprintf("A%v", rowStartIdx), time.Month(v.Month))
			f.SetCellInt(defaultSheet, fmt.Sprintf("B%v", rowStartIdx), v.Year)
			f.SetCellFloat(defaultSheet, fmt.Sprintf("C%v", rowStartIdx), v.Prediction, 2, 64)
			f.SetCellValue(defaultSheet, fmt.Sprintf("D%v", rowStartIdx), kWhType)

			rowStartIdx += 1
		}

		f.SetCellValue(defaultSheet, fmt.Sprintf("A%v", rowStartIdx), fmt.Sprintf("MAPE: %v", res.MAPE))
		rowStartIdx += 1
	}

	f.SetActiveSheet(index)
	// save spreadsheet by the given path
	if err := f.SaveAs("result.xlsx"); err != nil {
		return err
	}

	return nil
}

func export(res []*Result) error {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	var (
		defaultSheet = "Sheet1"
	)

	// create a new sheet
	index, err := f.NewSheet(defaultSheet)
	if err != nil {
		return err
	}

	// START CELL STYLE =========================================================================================== //

	centeredStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})

	boldYellowBgStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#ffff00"}, Pattern: 1},
	})

	// boldBlueBgStyle, _ := f.NewStyle(&excelize.Style{
	// 	Font: &excelize.Font{
	// 		Bold:  true,
	// 		Color: "#ffffff",
	// 	},
	// 	Fill: excelize.Fill{Type: "pattern", Color: []string{"#1135d4"}, Pattern: 1},
	// })

	boldCenteredStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})

	customNumFmt := "#,##0.##"
	thousandSeparatorNumberStyle, _ := f.NewStyle(&excelize.Style{CustomNumFmt: &customNumFmt})

	// END CELL STYLE =========================================================================================== //

	// START DEFINING ROW/COL SIZE =========================================================================================== //

	f.SetColWidth(defaultSheet, "A", "A", 5)
	f.SetColWidth(defaultSheet, "B", "B", 20)
	f.SetColWidth(defaultSheet, "D", "O", 20)

	// END DEFINING ROW/COL SIZE =========================================================================================== //

	// START PRINTING PREDICTIONS =========================================================================================== //

	rowStartIdx := 3 // starting point
	for i, v := range res {
		f.SetCellValue(defaultSheet, fmt.Sprintf("B%v", rowStartIdx+i), fmt.Sprintf("FORJA %v", v.MAPE))
		f.SetCellValue(defaultSheet, fmt.Sprintf("D%v", rowStartIdx+1+i), "Bulan")

		f.MergeCell(defaultSheet, fmt.Sprintf("D%v", rowStartIdx+1+i), fmt.Sprintf("O%v", rowStartIdx+1+i))
		f.MergeCell(defaultSheet, fmt.Sprintf("A%v", rowStartIdx+1+i), fmt.Sprintf("A%v", rowStartIdx+2+i))
		f.MergeCell(defaultSheet, fmt.Sprintf("B%v", rowStartIdx+1+i), fmt.Sprintf("B%v", rowStartIdx+2+i))
		f.MergeCell(defaultSheet, fmt.Sprintf("C%v", rowStartIdx+1+i), fmt.Sprintf("C%v", rowStartIdx+2+i))

		f.SetCellStyle(defaultSheet, fmt.Sprintf("A%v", rowStartIdx+1+i), fmt.Sprintf("O%v", rowStartIdx+2+i), boldCenteredStyle)
		f.SetCellStyle(defaultSheet, fmt.Sprintf("B%v", rowStartIdx+i), fmt.Sprintf("B%v", rowStartIdx+i), boldYellowBgStyle)

		no := 1
		valueStartRowIdx := rowStartIdx + 3
		for j, w := range v.Items {

			monthNameStartColIdx := 4 // D
			monthNameStartRowIdx := rowStartIdx + 2

			f.SetCellValue(defaultSheet, fmt.Sprintf("A%v", monthNameStartRowIdx+i), "No.")
			f.SetCellValue(defaultSheet, fmt.Sprintf("B%v", monthNameStartRowIdx+i), "Input")
			f.SetCellValue(defaultSheet, fmt.Sprintf("C%v", monthNameStartRowIdx+i), "Satuan")

			f.SetCellValue(defaultSheet, fmt.Sprintf("A%v", valueStartRowIdx+i), no)
			f.SetCellValue(defaultSheet, fmt.Sprintf("B%v", valueStartRowIdx+i), j)
			f.SetCellValue(defaultSheet, fmt.Sprintf("C%v", valueStartRowIdx+i), "kWh")

			f.SetCellStyle(defaultSheet, fmt.Sprintf("A%v", valueStartRowIdx+i), fmt.Sprintf("A%v", valueStartRowIdx+i), centeredStyle)
			f.SetCellStyle(defaultSheet, fmt.Sprintf("C%v", valueStartRowIdx+i), fmt.Sprintf("C%v", valueStartRowIdx+i), centeredStyle)

			for m := time.January; m <= time.December; m++ {
				f.SetCellValue(defaultSheet, fmt.Sprintf("%s%v", toChar(monthNameStartColIdx), monthNameStartRowIdx+i), m.String())
				f.SetCellValue(defaultSheet, fmt.Sprintf("%s%v", toChar(monthNameStartColIdx), valueStartRowIdx+i), math.Round(w.Prediction))

				f.SetCellStyle(defaultSheet, fmt.Sprintf("%s%v", toChar(monthNameStartColIdx), valueStartRowIdx+i), fmt.Sprintf("%s%v", toChar(monthNameStartColIdx), valueStartRowIdx+i), thousandSeparatorNumberStyle)

				monthNameStartColIdx += 1
			}

			valueStartRowIdx += 1

			no += 1

		}

		rowStartIdx += 14 // 14 is total length of populated rows for one year
	}

	// END PRINTING PREDICTIONS =========================================================================================== //

	f.SetActiveSheet(index)
	// save spreadsheet by the given path
	if err := f.SaveAs("result.xlsx"); err != nil {
		return err
	}

	return nil
}

func toChar(i int) string {
	return strings.Replace(string(rune('A'-1+i)), "'", "", -1)
}
