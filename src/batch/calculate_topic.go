package batch

import (
	"fmt"
	"log"
	//"log"
	//"testing"

	. "ml"
	. "model"
)

func Calculate_topic() (err error) {
	var langs []string = []string{
		"ja",
		"ko",
		"cn",
		"en",
	}

	var ttM TokenizedTweetsModel = TokenizedTweetsModel{}
	ttM.Init()

	var twTMs TweetTopicModel = TweetTopicModel{}

	//var result []uint32 = make([]uint32, 0, 4)
	//var words []string = make([]string, 0, 4)

	for _, lang := range langs {
		//get hash word map
		cur2, err := ttM.FindHashWordMapByLang(lang)

		if err != nil {
			return err
		}

		defer cur2.Close(ttM.GetCtx())
		var hashWordMap MongoHashWordMap
		for cur2.Next(ttM.GetCtx()) {
			err = cur2.Decode(&hashWordMap)

			if err != nil {
				return err
			}
		}

		cur, err := ttM.FindTfIdfsByLang(lang)

		if err != nil {
			return err
		}

		var results []MongoCalculatedTfIdfTokenizedTweets = make([]MongoCalculatedTfIdfTokenizedTweets, 0, 1000000)
		defer cur.Close(ttM.GetCtx())
		for cur.Next(ttM.GetCtx()) {
			var result MongoCalculatedTfIdfTokenizedTweets
			err = cur.Decode(&result)

			if err != nil {
				return err
			}

			results = append(results, result)
		}

		var calculatedTfIdfs []CalculatedTfIdfForTopicModel = make([]CalculatedTfIdfForTopicModel, 0, 1000000)
		for _, result := range results {
			for _, tfIdfData := range result.TfIdfs {
				word := hashWordMap.Values.Values[fmt.Sprintf("%d", tfIdfData.Hash)]

				if len(word) > 3 {
					var calculatedTfIdf CalculatedTfIdfForTopicModel = CalculatedTfIdfForTopicModel{
						Hash:  tfIdfData.Hash,
						TfIdf: tfIdfData.TfIdf,
					}
					calculatedTfIdfs = append(calculatedTfIdfs, calculatedTfIdf)
				}
			}
		}

		//calculate topic
		_, pickedUpScoreList, err := CalculateTopic(&calculatedTfIdfs)

		//twTMs.Begin()

		err = twTMs.DeleteTopic(lang)

		if err != nil {
			return err
		}

		log.Printf("lang:%s pickedUpScoreList:%d", lang, len(pickedUpScoreList))
		for hashNum, score := range pickedUpScoreList {
			word := hashWordMap.Values.Values[fmt.Sprintf("%d", hashNum)]
			twTMs.Topic = word
			twTMs.Score = score
			twTMs.Lang = lang
			twTMs.InsertTopic()
		}

		//twTMs.Commit()
		//words = append(words, hashWordMap.Values.Values[hashNum])
		//log.Println(hashWordMap.Values.Values[fmt.Sprintf("%d", hashNum)])

		/**
		for _, hashNum := range hashNumList {
			log.Printf("%s/%f", hashWordMap.Values.Values[fmt.Sprintf("%d", hashNum)], pickedUpScoreList[hashNum])
		}
		words = append(words, "hoge")
		**/
	}

	return
}
