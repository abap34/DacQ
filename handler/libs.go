package handler

import (
	"encoding/csv"
	"fmt"
	"github.com/pkg/errors"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
)

var (
	csvFiles  = make(map[string]*CSVFile)
	templates = template.Must(template.ParseFiles("templates/index.html"))
	rankings  []Ranking
)

func conguraturation(w http.ResponseWriter, r *http.Request, user string, loss float64, isBest bool) {
	fmt.Println("Conguraturation!")
	redirectURL := "/result?user=" + url.QueryEscape(user) + "&loss=" + strconv.FormatFloat(loss, 'f', -1, 64) + "&isBest=" + strconv.FormatBool(isBest)
	fmt.Println("redirectURL: ", redirectURL)
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func sortRankings() {
	sort.Slice(rankings, func(i, j int) bool {
		return rankings[i].Loss < rankings[j].Loss
	})

	for i := range rankings {
		rankings[i].Rank = i + 1
	}
}

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

func getSessionUser(r *http.Request) string {
	user := r.Header.Get("X-Forwarded-User")
	if user == "" {
		user = "名無しのエンジニア"
	}

	return user
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
