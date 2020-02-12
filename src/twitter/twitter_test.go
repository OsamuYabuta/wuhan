package twitter

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"testing"
	"time"
)

/**
twitter accesstoken  1222513666467090432-qitxN1tKJ6SxyfkpTa0FAxQvQNhhmO
twitter accesstoken secret Rs60rro7zZDGtz8Jgyhnq0E1m45zd6gtaCa0eAFZTUqcm
api key BUE5F9Vi0HoRUCBQnwKrUwdxG
api secret key PvYtpzU6me6YUazMdCZF0ooPU6n4SnzN22XrUm524AdZOYols8
**/

func TestSearch(t *testing.T) {

	//var langs []string = []string{"ja","en","ch","ko"}
	//var keywords []string = []string{"武漢","wuhan","武汉","우한"}
	var langs []string = []string{"ko", "zh-cn"}
	var keywords []string = []string{"우한", "武汉"}

	for i, l := range langs {
		var tw Twitter = Twitter{}
		tw.Init()
		tw.SetApiKey("BUE5F9Vi0HoRUCBQnwKrUwdxG")
		tw.SetApiSecretKey("PvYtpzU6me6YUazMdCZF0ooPU6n4SnzN22XrUm524AdZOYols8")
		tw.SetOauthToken("1222513666467090432-qitxN1tKJ6SxyfkpTa0FAxQvQNhhmO")
		tw.SetOauthTokenSecret("Rs60rro7zZDGtz8Jgyhnq0E1m45zd6gtaCa0eAFZTUqcm")
		ret, err := tw.SearchForKeywordInTweets(keywords[i], l)

		t.Log(len(ret.Statuses))

		if err != nil {
			t.Fatal(err.Error())
		}

		text := ""
		for _, v := range ret.Statuses {

			text += fmt.Sprintf("%d/%s/%s/%s\n", v.Id, v.CreatedAt, v.Text, v.User.Username)
		}

		ioutil.WriteFile("./"+l+"_result.txt", []byte(text), 0777)

		//time.Sleep(time.Millisecond * 5000)

		var i int = 0
		for {
			ret2, err := tw.SearchNext()

			if err != nil {
				t.Fatal(err.Error())
			}

			text = ""
			for _, v := range ret2.Statuses {
				text += fmt.Sprintf("%d/%s/%s\n", v.Id, v.CreatedAt, v.Text)
			}

			ioutil.WriteFile("./"+l+"_result"+strconv.Itoa(i)+".txt", []byte(text), 0777)
			time.Sleep(time.Millisecond * 1000)

			i++
			if i > 1 {
				break
			}
		}
	}
}
