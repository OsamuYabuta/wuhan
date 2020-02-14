package batch

import (
	"log"
	"sync"
	//"testing"
	"utils"

	. "model"
	. "twitter"
)

func Collect_tweet() (err error) {
	var langs []string = []string{
		"ja",
		"zh-cn",
		"ko",
		"en",
	}
	var keywords []string = []string{
		"武漢",
		"武汉",
		"우한",
		"wuhan",
	}

	var wg sync.WaitGroup
	for i, lang := range langs {
		keyword := keywords[i]
		wg.Add(1)
		go sub(keyword, lang, &wg)
	}

	wg.Wait()

	return nil
}

func sub(keyword string, lang string, wg *sync.WaitGroup) {
	defer wg.Done()

	var tw Twitter = Twitter{}
	tw.Init()
	tw.SetApiKey("BUE5F9Vi0HoRUCBQnwKrUwdxG")
	tw.SetApiSecretKey("PvYtpzU6me6YUazMdCZF0ooPU6n4SnzN22XrUm524AdZOYols8")
	tw.SetOauthToken("1222513666467090432-qitxN1tKJ6SxyfkpTa0FAxQvQNhhmO")
	tw.SetOauthTokenSecret("Rs60rro7zZDGtz8Jgyhnq0E1m45zd6gtaCa0eAFZTUqcm")

	var twMs TweetModel = TweetModel{}
	sinceId, err := twMs.FindSinceId(lang)

	if err != nil {
		panic(err.Error())
	}

	tweets, err := tw.SearchForKeywordInTweets(keyword, lang, sinceId)

	if err != nil {
		panic(err.Error())
	}

	for _, tweet := range tweets.Statuses {
		twMs.Id = tweet.Id
		twMs.Username = tweet.User.Username
		twMs.Lang = tweet.Lang
		twMs.Tweet = tweet.Text
		twMs.CreatedAt = utils.ParseTweetedTime(tweet.CreatedAt)
		twMs.Screenname = tweet.User.Screenname

		//if it error , execute panic automatically
		twMs.Insert()
	}

	for {
		//ioutil.WriteFile("./limit.txt", []byte(tw.CurrentXRateLimitRemaining), 0777)
		log.Printf("rate limit remaining %s", tw.CurrentXRateLimitRemaining)
		if tw.CurrentXRateLimitRemaining == "0" {
			log.Printf("exceed rate limit of twitter search api.")
			break
		}

		tweets, err := tw.SearchNext()

		if err != nil {
			panic(err.Error())
		}

		for _, tweet := range tweets.Statuses {
			twMs.Id = tweet.Id
			twMs.Username = tweet.User.Username
			twMs.Lang = tweet.Lang
			twMs.Tweet = tweet.Text
			twMs.CreatedAt = utils.ParseTweetedTime(tweet.CreatedAt)
			twMs.Screenname = tweet.User.Screenname

			//if it error , execute panic automatically
			twMs.Insert()
		}
	}
}
