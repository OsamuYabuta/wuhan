package batch

import (
	//"testing"
	"time"
	//"time"
	"tokenizer"

	//"go.mongodb.org/mongo-driver/bson/primitive"
	. "ml"
	. "model"
	. "sync"
)

func Calculate_tfidf() (err error) {
	var ttM TokenizedTweetsModel = TokenizedTweetsModel{}
	ttM.Init()

	var langs []string = []string{
		"ja",
		"ko",
		"cn",
		"en",
	}

	var wg WaitGroup
	for _, lang := range langs {
		wg.Add(1)
		go sub3(&ttM, lang, &wg)
	}

	wg.Wait()

	return
}

func sub3(ttM *TokenizedTweetsModel, lang string, wg *WaitGroup) {
	defer wg.Done()

	cur, err := ttM.FindByLang(lang, time.Now().AddDate(0, 0, -1))

	if err != nil {
		panic(err.Error())
	}

	var results []MongoTokenizedTweets = make([]MongoTokenizedTweets, 0, 1000000)
	//ct := 0
	defer cur.Close(ttM.GetCtx())
	for cur.Next(ttM.GetCtx()) {
		var result MongoTokenizedTweets
		err = cur.Decode(&result)
		if err != nil {
			panic(err.Error())
		}

		results = append(results, result)
	}

	var matrix Matrix = Matrix{
		Rows: make([]Vector, len(results)),
	}

	var tweetIds []int64 = make([]int64, len(results))
	//ct = 0
	for i, result := range results {
		matrix.Rows[i].WordTf = make([]WordTf, len(result.Tokens))
		for i2, token := range result.Tokens {
			tToken := tokenizer.Token{
				Keyword: token.Keyword,
				Tf:      0,
				Tag:     token.Tag,
			}
			if tokenizer.IsNoun(lang, tToken) == true {
				matrix.Rows[i].WordTf[i2] = WordTf{
					Word: token.Keyword,
					Tf:   token.Tf,
				}
			}
		}
		tweetIds[i] = result.TweetId
	}

	//calculate tf-idf
	res, hashWordMap, err := CalculateTfIdf(matrix)

	if err != nil {
		panic(err.Error())
	}

	var wg2 WaitGroup
	for i, row := range res.Rows {
		//get object id of mongo docuemnt
		TweetId := tweetIds[i]
		wg2.Add(1)
		go ttM.InsertTfIdfTweets(TweetId, lang, row, &hashWordMap, &wg2)
	}

	err = ttM.InsertHashWordMap(lang, &hashWordMap)

	if err != nil {
		panic(err.Error())
	}

	wg2.Wait()
}
