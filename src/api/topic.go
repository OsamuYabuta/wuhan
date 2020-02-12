package api

import (
	"encoding/json"
	"model"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

type ApiTopicResultJson struct {
	Lang string  `json:"lang"`
	Data []Topic `json:"topics"`
}

type Topic struct {
	Id             int64   `json:"id"`
	Topic          string  `json:"topic"`
	Score          float64 `json:"score"`
	Lang           string  `json:"lang"`
	CalculatedDate string  `json:"cdate"`
}

func Api_topic(writter http.ResponseWriter, request *http.Request, param httprouter.Params) {
	lang := param.ByName("lang")

	var ttM model.TweetTopicModel = model.TweetTopicModel{}
	var targetDate = time.Now().AddDate(0, 0, -1)
	result, err := ttM.FindByLang(lang, targetDate)

	if err != nil {
		writter.Header().Set("Access-Control-Allow-Origin", "*")
		writter.Write([]byte("{error:\"" + err.Error() + "\"}"))
		return
	}

	var resJson ApiTopicResultJson = ApiTopicResultJson{
		Lang: lang,
		Data: make([]Topic, 0, len(result)),
	}

	for _, rec := range result {
		resJson.Data = append(resJson.Data, Topic{
			Id:             rec.Id,
			Lang:           rec.Lang,
			Topic:          rec.Topic,
			Score:          rec.Score,
			CalculatedDate: rec.CalculatedDate,
		})
	}

	resJsonByte, err := json.MarshalIndent(&resJson, "", "\t\t")

	if err != nil {
		writter.Header().Set("Access-Control-Allow-Origin", "*")
		writter.Write([]byte("{error:\"" + err.Error() + "\"}"))
		return
	}

	writter.Header().Set("Access-Control-Allow-Origin", "*")
	writter.Write(resJsonByte)
}
