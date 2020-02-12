package oauth

import (
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

func TestOauth(t *testing.T) {
	var oa OAuth = OAuth{}
	oa.Init()
	oa.SetApiBaseUrl("https://api.twitter.com/oauth/request_token")
	//oa.SetApiBaseUrl("https://api.twitter.com/1.1/statuses/update.json")
	oa.SetRequestMethod("POST")
	//oa.SetParameter("include_entities", "true")
	//oa.SetParameter("status", "Hello Ladies + Gentlemen, a signed OAuth request!")
	//oa.SetApiBaseUrl("http://o-freedom.info")
	oa.SetCallback("http://o-freedom.info/wuhan/oauth/callback")
	//oa.SetCallback("") //http://myapp.com:3005/twitter/process_callback")
	oa.SetConsumerKey("BUE5F9Vi0HoRUCBQnwKrUwdxG")
	//oa.SetConsumerKey("xvz1evFS4wEEPTGEFPHBog")
	oa.SetConsumerSecretKey("PvYtpzU6me6YUazMdCZF0ooPU6n4SnzN22XrUm524AdZOYols8")
	//oa.SetConsumerSecretKey("kAcSOqF21Fu85e7zjz7ZN2U4ZRhfV3WpwPAoE3Z7kBw")
	//oa.SetOauthToken("370773112-GmHxMAgYyLbNEtIKZeRNFsMKPR9EyMZeS9weJAEb")
	//oa.SetOauthTokenSecret("LswwdoUaIvS8ltyTt5jkRh4J50vUPVVHtR2YPi5kE")

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	//postbody := strings.NewReader(fmt.Sprintf("%s=%s", OAUTH_CALLBACK_KEY, url.QueryEscape(oa.GetCallback())))
	req, err := http.NewRequest("POST", "https://api.twitter.com/oauth/request_token", nil)
	//req, err := http.NewRequest("POST", "http://o-freedom.info", nil)
	if err != nil {
		t.Fatal(err.Error())
	}
	oa.SetAuthorizationHeader(req)
	res, err := client.Do(req)

	if err != nil {
		t.Fatal(err.Error())
	}

	os.Create("./result.txt")
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		t.Fatal(err.Error())
	}

	ioutil.WriteFile("./result.txt", body, 0777)
}
