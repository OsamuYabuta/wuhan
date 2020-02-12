package batch

import (
	"log"
	//"testing"

	. "model"
	. "sync"
	. "tokenizer"
)

func Tokenize_tweets() (err error) {
	var twMs TweetModel = TweetModel{}
	var langs []string = []string{
		"ja",
		"ko",
		"zh",
		"en",
	}

	var ttM TokenizedTweetsModel = TokenizedTweetsModel{}
	ttM.Init()

	var wg WaitGroup
	for _, lang := range langs {
		result, err := twMs.FindByLang(lang)

		if err != nil {
			return err
		}

		wg.Add(1)
		go sub2(&ttM, result, lang, &wg)
	}

	wg.Wait()

	return
}

func sub2(ttM *TokenizedTweetsModel, result []TweetModel, lang string, wg *WaitGroup) {
	defer wg.Done()

	if lang == "zh" {
		lang = "cn"
	}

	var tokenizer Tokenizer = Tokenizer{}

	log.Printf("lang:%s count:%d", lang, len(result))

	var wg2 WaitGroup
	for _, tweet := range result {

		if has, err := ttM.HasTokenizedTweet(tweet.Id); has == true {
			continue
		} else if err != nil {
			panic(err.Error())
		}

		result2, err := tokenizer.Tokenize(tweet.Tweet, lang)

		if err != nil {
			panic(err.Error())
		}

		Clearning(&result2)

		wg2.Add(1)
		go ttM.InsertTokenizedTweets(tweet.Id, lang, tweet.CreatedAt, &result2, &wg2)
		//err := ttM.InsertTokenizedTweets(&tokenizedTweets)
	}

	wg2.Wait()

	log.Printf("lang:%s is finished to process", lang)
}
