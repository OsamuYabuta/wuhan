package model

import (
	"fmt"
	"sync"
	"testing"
	"tokenizer"

	"go.mongodb.org/mongo-driver/bson"
)

func TestTweetsSql(t *testing.T) {
	var twm TweetModel = TweetModel{}

	for i := 0; i < 10; i++ {
		twm.Id = int64(i)
		twm.Username = fmt.Sprintf("user%d", i)
		twm.Tweet = fmt.Sprintf("これはツイート%dです", i)
		twm.Lang = "ja"
		twm.CreatedAt = "2020-02-04 00:00:00"

		twm.Insert()
	}
}

func TestFindByLang(t *testing.T) {
	var twm TweetModel = TweetModel{}

	result, err := twm.FindByLang("ja")

	if err != nil {
		t.Fatal(err.Error())
	}

	t.Log(fmt.Sprintf("tweets size:%d", len(result)))

	for _, tweet := range result {
		t.Log(tweet.Username + "/" + tweet.Tweet)
	}

	t.Fatal("ended test")
}

func TestMongo(t *testing.T) {
	var ttM TokenizedTweetsModel = TokenizedTweetsModel{}
	ttM.Init()

	var tokens tokenizer.Tokens = tokenizer.Tokens{
		Values: []tokenizer.Token{
			tokenizer.Token{
				Keyword: "キーワード1",
				Tag:     "名詞",
				Tf:      1,
			},
			tokenizer.Token{
				Keyword: "キーワード2",
				Tag:     "動詩",
				Tf:      2,
			},
		},
	}

	var wg sync.WaitGroup
	wg.Add(1)
	ttM.InsertTokenizedTweets(100, "ja", &tokens, &wg)
	wg.Wait()
}

func TestMongFindByLang(t *testing.T) {
	var ttM TokenizedTweetsModel = TokenizedTweetsModel{}
	ttM.Init()

	cur, err := ttM.FindByLang("ja")

	if err != nil {
		t.Fatal(err.Error())
	}

	var results []bson.M
	for cur.Next(ttM.GetCtx()) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			t.Fatal(err.Error())
		}
		results = append(results, result)
	}

	t.Log(results)

	t.Fatal("ended")
}

func TestFindSinceId(t *testing.T) {
	var twM TweetModel = TweetModel{}

	sinceId, err := twM.FindSinceId("ja")

	if err != nil {
		panic(err.Error())
	}
}
