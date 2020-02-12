package twitter

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"oauth"
	"strings"
)

const (
	SEARCH_API_ENDPOINT = "https://api.twitter.com/1.1/search/tweets.json"
)

type Tweets struct {
	Statuses       []Tweet        `json:"statuses"`
	SearchMetaData SearchMetaData `json:"search_metadata"`
}

type Tweet struct {
	CreatedAt string      `json:"created_at"`
	Id        int64       `json:"id"`
	Text      string      `json:"text"`
	Lang      string      `json:"lang"`
	User      TwitterUser `json:"user"`
}

type TwitterUser struct {
	Username   string `json:"name"`
	Screenname string `json:"screen_name"`
}

type SearchMetaData struct {
	MaxIdStr    string `json:"max_id_str"`
	NextResults string `json:"next_results"`
	Query       string `json:"query"`
	SinceId     int64  `json:"since_id"`
}

type TwitterResponseErrors struct {
	Errors []TwitterResponseError `json:"errors"`
}

type TwitterResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Twitter struct {
	OA                         oauth.OAuth
	KeyValues                  map[string]string
	PrevSearchMetaData         SearchMetaData
	PrevTweetIds               []int64
	CurrentXRateLimitLimit     string
	CurrentXRateLimitRemaining string
	CurrentXRateLimitReset     string
}

func (tw *Twitter) Init() {
	tw.OA = oauth.OAuth{}
	tw.OA.Init()
	tw.KeyValues = make(map[string]string, 20)
}

func (tw *Twitter) Clean() {
	tw.OA.Clear()
	tw.KeyValues = make(map[string]string, 20)
}

func (tw *Twitter) SetApiKey(apiKey string) {
	tw.OA.SetConsumerKey(apiKey)
}

func (tw *Twitter) SetApiSecretKey(apiSecretKey string) {
	tw.OA.SetConsumerSecretKey(apiSecretKey)
}

func (tw *Twitter) SetOauthToken(oauthToken string) {
	tw.OA.SetOauthToken(oauthToken)
}

func (tw *Twitter) SetOauthTokenSecret(oauthSecret string) {
	tw.OA.SetOauthTokenSecret(oauthSecret)
}

func (tw *Twitter) SetParameter(key string, value string) {
	tw.KeyValues[key] = value
	tw.OA.SetParameter(key, value)
}

func (tw *Twitter) buildQueryString(keyValues *map[string]string) string {
	var queryParameterSlice []string
	for k, v := range *keyValues {
		if k != "" && v != "" {
			queryParameterSlice = append(queryParameterSlice, fmt.Sprintf("%s=%s", k, url.QueryEscape(v)))
		}
	}

	return strings.Join(queryParameterSlice, "&")
}

func (tw *Twitter) SearchForKeywordInTweets(keyword string, lang string, sinceId int64) (result Tweets, err error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	//prepare parameter
	tw.OA.SetApiBaseUrl(SEARCH_API_ENDPOINT)
	tw.OA.SetCallback("")
	tw.OA.SetRequestMethod("GET")
	tw.SetParameter("q", keyword)
	tw.SetParameter("result_type", "mixed")
	tw.SetParameter("lang", lang)
	tw.SetParameter("count", "100")
	if sinceId > 0 {
		tw.SetParameter("since_id", fmt.Sprintf("%d", sinceId))
	}
	queryString := tw.buildQueryString(&tw.KeyValues)

	req, err := http.NewRequest("GET", fmt.Sprintf("%s?%s", SEARCH_API_ENDPOINT, queryString), nil)
	if err != nil {
		return result, err
	}

	//set authorization header for oauth
	err = tw.OA.SetAuthorizationHeader(req)
	if err != nil {
		return result, err
	}

	res, err := client.Do(req)
	if err != nil {
		return result, err
	}

	defer res.Body.Close()

	rawJsonResult, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return result, err
	}

	//ioutil.WriteFile("./rawjson.txt", []byte(rawJsonResult), 0777)

	errorResult := TwitterResponseErrors{}
	json.Unmarshal(rawJsonResult, &errorResult)
	if len(errorResult.Errors) > 0 {
		return result, errors.New(fmt.Sprintf("%d/%s", errorResult.Errors[0].Code, errorResult.Errors[0].Message))
	}

	tw.CurrentXRateLimitLimit = res.Header.Get("x-rate-limit-limit")
	tw.CurrentXRateLimitRemaining = res.Header.Get("x-rate-limit-remaining")
	tw.CurrentXRateLimitReset = res.Header.Get("x-rate-limit-reset")

	json.Unmarshal(rawJsonResult, &result)
	tw.PrevSearchMetaData = result.SearchMetaData

	tw.Clean()

	return result, nil
}

func (tw *Twitter) SearchNext() (result Tweets, err error) {
	prevSearchMetaData := tw.PrevSearchMetaData
	nextResults := prevSearchMetaData.NextResults

	//sinceId := tw.PrevTweetIds[0] - 1

	//log.Println(nextResults)
	var nextResultKeyValues map[string]string = make(map[string]string)
	var nextResultsSlice []string = strings.Split(strings.ReplaceAll(nextResults, "?", ""), "&")
	for _, v := range nextResultsSlice {
		keyValue := strings.Split(v, "=")
		if len(keyValue) != 2 {
			return result, errors.New("invalid parse next results in search meta data.")
		}
		key := keyValue[0]
		value, err := url.QueryUnescape(keyValue[1])
		if err != nil {
			return result, err
		}
		nextResultKeyValues[key] = value
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	//prepare parameter
	tw.OA.SetApiBaseUrl(SEARCH_API_ENDPOINT)
	tw.OA.SetCallback("")
	tw.OA.SetRequestMethod("GET")
	//tw.SetParameter("since_id", fmt.Sprintf("%d", sinceId))

	for k, v := range nextResultKeyValues {
		if k == "include_entities" {
			switch v {
			case "1":
				v = "true"
			default:
				v = "false"
			}
		}
		tw.SetParameter(k, v)
	}

	tw.SetParameter("count", "100")

	queryString := tw.buildQueryString(&tw.KeyValues)

	req, err := http.NewRequest("GET", fmt.Sprintf("%s?%s", SEARCH_API_ENDPOINT, queryString), nil)
	if err != nil {
		return result, err
	}

	//set authorization header for oauth
	err = tw.OA.SetAuthorizationHeader(req)
	if err != nil {
		return result, err
	}

	res, err := client.Do(req)
	if err != nil {
		return result, err
	}

	rawJsonResult, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return result, err
	}

	errorResult := TwitterResponseErrors{}
	json.Unmarshal(rawJsonResult, &errorResult)
	if len(errorResult.Errors) > 0 {
		return result, errors.New(fmt.Sprintf("Code:%d/Message:%s", errorResult.Errors[0].Code, errorResult.Errors[0].Message))
	}

	tw.CurrentXRateLimitLimit = res.Header.Get("x-rate-limit-limit")
	tw.CurrentXRateLimitRemaining = res.Header.Get("x-rate-limit-remaining")
	tw.CurrentXRateLimitReset = res.Header.Get("x-rate-limit-reset")

	json.Unmarshal(rawJsonResult, &result)
	tw.PrevSearchMetaData = result.SearchMetaData
	tw.Clean()

	return result, nil
}
