package batch

import (
	//	"testing"

	. "ml"
	. "model"
)

func Calculate_User() (err error) {
	var langs []string = []string{
		"ja",
		"ko",
		"cn",
		"en",
	}

	var twM TweetModel = TweetModel{}
	var ttfTwM TokenizedTweetsModel = TokenizedTweetsModel{}
	var twPuMs TweetPickupUsersModel = TweetPickupUsersModel{}
	ttfTwM.Init()
	var userNames map[string]bool = make(map[string]bool)
	for _, lang := range langs {
		if lang == "cn" {
			lang = "zh"
		}
		result, err := twM.FindByLang(lang)

		if err != nil {
			return err
		}

		var tokens []TokenizedTweetsForPickupUser = make([]TokenizedTweetsForPickupUser, 0, 1000000)
		for _, tweet := range result {
			userNames[tweet.Screenname] = true
			tweetId := tweet.Id

			cur, err := ttfTwM.FindByTweetId(tweetId)

			if err != nil {
				return err
			}

			defer cur.Close(ttfTwM.GetCtx())
			for cur.Next(ttfTwM.GetCtx()) {
				var result2 MongoTokenizedTweets
				err = cur.Decode(&result2)

				if err != nil {
					return err
				}

				for _, token := range result2.Tokens {
					if len(token.Keyword) > 3 {
						tokens = append(tokens, TokenizedTweetsForPickupUser{
							Keyword: token.Keyword,
							Tag:     token.Tag,
							Tf:      token.Tf,
						})
					}
				}
			}
		}
		pickedUpUsers := CalculatePickupUser(&userNames, &tokens)

		err = twPuMs.DeletePickupUsers(lang)

		if err != nil {
			return err
		}

		for _, keyValue := range pickedUpUsers {
			twPuMs.Lang = lang
			//log.Printf("%d", keyValue.Value)
			twPuMs.Score = float32(keyValue.Value) * 1.0
			twPuMs.Screenname = keyValue.Key
			twPuMs.InsertPickedUpUser()
		}

	}

	return
}
