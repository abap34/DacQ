package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func PostUploadCSV(w http.ResponseWriter, r *http.Request) {
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

	for i, ranking := range rankings {
		if ranking.User == user {
			if csvFile.Loss < ranking.Loss {
				rankings[i].Loss = csvFile.Loss
				conguraturation(w, r, user, csvFile.Loss, true)
				fmt.Println("Conguraturation!")
				return
			} else {
				fmt.Println("not Conguraturation!")
				conguraturation(w, r, user, csvFile.Loss, false)
				return
			}
		}
	}

	rankings = append(rankings, Ranking{User: csvFile.User, Loss: csvFile.Loss})

	http.Redirect(w, r, "/", http.StatusFound)
}
