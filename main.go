package main

import (
	"embed"
	_ "embed"
	"encoding/csv"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// CSVFile represents the CSV file uploaded by the user.
type CSVFile struct {
	User string
	File *os.File
	Loss float64
}

// Ranking represents the user ranking based on the loss.
type Ranking struct {
	User string
	Loss float64
	Rank int
}

var (
	csvFiles = make(map[string]*CSVFile)
	rankings []Ranking

	//go:embed templates/index.html
	indexHtml embed.FS
	templates = template.Must(template.ParseFS(indexHtml, "templates/index.html"))
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/upload-csv", uploadCSVHandler).Methods("POST")

	staticHandler := http.FileServer(http.Dir("static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", staticHandler))

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8000", r))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	sortRankings()

	data := struct {
		User     string
		Rankings []Ranking
	}{
		User:     getSessionUser(r),
		Rankings: rankings,
	}

	err := templates.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func uploadCSVHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20) // Limit file size to 32MB
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := r.Header.Get("X-Forwarded-User")
	if user == "" {
		user = r.Header.Get("X-Showcase-User")
	}

	file, handler, err := r.FormFile("csv_file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	osFile, err := os.Create(handler.Filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer osFile.Close()

	_, err = io.Copy(osFile, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	csvFile := &CSVFile{User: user, File: osFile}
	csvFiles[user] = csvFile

	err = calculateLoss(csvFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rankings = append(rankings, Ranking{User: csvFile.User, Loss: csvFile.Loss})
	http.Redirect(w, r, "/", http.StatusFound)
}

func calculateLoss(csvFile *CSVFile) error {
	// ファイルの読み込み
	labelFile, err := os.Open("label.csv")
	if err != nil {
		return errors.Wrap(err, "failed to open label file")
	}
	defer labelFile.Close()

	// labelファイルのデータを読み取る
	labelReader := csv.NewReader(labelFile)
	labelRecords, err := labelReader.ReadAll()
	if err != nil {
		return errors.Wrap(err, "failed to read label file")
	}

	// 投稿されたCSVファイルのデータを読み取る
	csvFile.File.Seek(0, 0)
	csvReader := csv.NewReader(csvFile.File)
	csvRecords, err := csvReader.ReadAll()
	if err != nil {
		return errors.Wrap(err, "failed to read CSV file")
	}

	// ラベルと投稿データの比較と誤差の計算
	if len(labelRecords) != len(csvRecords) {
		return errors.New("label file and CSV file have different number of records")
	}

	actualValues := make([]float64, len(labelRecords)-1)
	predictedValues := make([]float64, len(csvRecords)-1)

	for i := 1; i < len(labelRecords); i++ {
		label, err := strconv.ParseFloat(labelRecords[i][0], 64)
		if err != nil {
			return errors.Wrapf(err, "failed to parse label value in row %d", i+1)
		}
		actualValues[i-1] = label

		predicted, err := strconv.ParseFloat(csvRecords[i][0], 64)
		if err != nil {
			return errors.Wrapf(err, "failed to parse predicted value in row %d", i+1)
		}
		predictedValues[i-1] = predicted
	}

	csvFile.Loss = meanSquaredError(actualValues, predictedValues)
	return nil
}

func meanSquaredError(actual, predicted []float64) float64 {
	if len(actual) != len(predicted) {
		return 9999999999999.999999999999
	}

	sum := 0.0
	for i := 0; i < len(actual); i++ {
		diff := actual[i] - predicted[i]
		sum += diff * diff
	}

	return sum / float64(len(actual))
}

func sortRankings() {
	sort.Slice(rankings, func(i, j int) bool {
		return rankings[i].Loss < rankings[j].Loss
	})

	for i := range rankings {
		rankings[i].Rank = i + 1
	}
}

func getSessionUser(r *http.Request) string {
	user := r.Header.Get("X-Forwarded-User")
	if user == "" {
		user = "名無しのエンジニア"
	}

	return user
}
