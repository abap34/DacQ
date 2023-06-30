package handler

import (
	"github.com/dacq/model"
	"net/http"
)

func GetHome(w http.ResponseWriter, r *http.Request) {
	scores, err := model.GetRankings()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rankings := make([]Ranking, len(scores))
	for i, score := range scores {
		rankings[i] = Ranking{
			User: score.User,
			Loss: score.Loss,
			Rank: i + 1,
		}
	}

	data := struct {
		User     string
		Rankings []Ranking
	}{
		User:     getSessionUser(r),
		Rankings: rankings,
	}

	err = templates.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
