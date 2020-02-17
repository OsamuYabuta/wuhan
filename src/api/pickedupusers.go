package api

import (
	"encoding/json"
	"model"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

type ApiPickedUpUsersResult struct {
	Lang string         `json:"lang"`
	Data []PickedUpUser `json:"pickedupusers"`
}

type PickedUpUser struct {
	Id             int64  `json:"id"`
	Lang           string `json:"lang"`
	Screenname     string `json:"screen_name"`
	Score          int    `json:"score"`
	CalculatedDate string `json:"cdate"`
}

func Api_Pickedupusers(writter http.ResponseWriter, request *http.Request, param httprouter.Params) {
	lang := param.ByName("lang")

	var tpkuM model.TweetPickupUsersModel = model.TweetPickupUsersModel{}
	var targetDate = time.Now()
	result, err := tpkuM.FindByLang(lang, targetDate)

	if err != nil {
		writter.Header().Set("Access-Control-Allow-Origin", "*")
		writter.Write([]byte("{error:\"" + err.Error() + "\"}"))
		return
	}

	var resJson ApiPickedUpUsersResult = ApiPickedUpUsersResult{
		Lang: lang,
		Data: make([]PickedUpUser, 0, len(result)),
	}

	for _, rec := range result {
		resJson.Data = append(resJson.Data, PickedUpUser{
			Id:             rec.Id,
			Lang:           rec.Lang,
			Screenname:     rec.Screenname,
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
