package handler

import (
	"math/rand"
	"net/http"
	"strconv"

	"github.com/dacq/model"
)

func GetHome(w http.ResponseWriter, r *http.Request) {
	scores, err := model.GetRankings()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rankings := make([]Ranking, len(scores))
	for i, score := range scores {
		user := score.User
		submitcount, err := model.CountSubmitLogByUser(user)
		lastsub, err := model.GetLastSubmitTime(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rankings[i] = Ranking{
			User:        score.User,
			Loss:        score.Loss,
			Rank:        i + 1,
			SubmitCount: submitcount,
			Last:        lastsub,
		}
	}

	messages := []string{
		"Kaggle班では新規班員を募集中! 詳しくは @abap34 または #team/kaggle まで！",
		"Welcome " + getSessionUser(r) + "!",
	}

	news, err := model.GetNews(5)

	for _, n := range news {
		news_text := "【ニュース】" + n.User + "がベストスコアを更新しました！ スコア:" + strconv.FormatFloat(n.Score, 'f', 6, 64)
		messages = append(messages, news_text)
	}

	// ランダムに並べ替え
	rand.Shuffle(len(messages), func(i, j int) {
		messages[i], messages[j] = messages[j], messages[i]
	})

	data := struct {
		User     string
		Rankings []Ranking
		Messages []string
	}{
		User:     getSessionUser(r),
		Rankings: rankings,
		Messages: messages,
	}

	err = templates.ExecuteTemplate(w, "index.html", data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
