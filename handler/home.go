package handler

import "net/http"

func GetHome(w http.ResponseWriter, r *http.Request) {
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
