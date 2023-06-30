package handler

import (
	"html/template"
	"net/http"
	"strconv"
)

func GetResult(w http.ResponseWriter, r *http.Request) {
	user := r.URL.Query().Get("user")
	lossStr := r.URL.Query().Get("loss")
	isBestStr := r.URL.Query().Get("isBest")
	isBest, err := strconv.ParseBool(isBestStr)
	if err != nil {
		http.Error(w, "Invalid isBest value", http.StatusBadRequest)
		return
	}
	loss, err := strconv.ParseFloat(lossStr, 64)
	if err != nil {
		http.Error(w, "Invalid loss value", http.StatusBadRequest)
		return
	}

	tmpl, err := template.ParseFiles("templates/result.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		User   string
		Loss   float64
		IsBest bool
	}{
		User:   user,
		Loss:   loss,
		IsBest: isBest,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
