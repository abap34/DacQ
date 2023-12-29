package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/dacq/model"
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
	defer func() {
		//delete
		err := os.Remove(osFile.Name())
		if err != nil {
			fmt.Println("failed to delete file: ", err.Error())
		}
		err = osFile.Close()
		if err != nil {
			fmt.Println("failed to close file: ", err.Error())
		}
	}()

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

	score, err := model.GetScoreByUser(user)

	if err != nil {
		if err.Error() != "record not found" {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		conguraturation(w, r, user, csvFile.Loss, true)
		fmt.Println("Conguraturation!")
		score := model.Score{User: user, Loss: csvFile.Loss}
		err = model.CreateScore(score)

		return
	}

	is_best := score.Loss > csvFile.Loss

	// イベントログをアップデート
	err = model.CreateSubmitLog(model.SubmitLog{
		User:   user,
		Time:   time.Now().Format("2004-03-04 10:01:01"),
		IsBest: is_best,
		Score:  csvFile.Loss,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if is_best {
		err = model.UpdateScore(score.User, csvFile.Loss)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		conguraturation(w, r, user, csvFile.Loss, true)
		fmt.Println("Conguraturation!")
	} else {
		fmt.Println("not Conguraturation!")
		conguraturation(w, r, user, csvFile.Loss, false)
	}

	return

}
