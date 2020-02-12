package tokenizer

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
	"unsafe"
)

type Tokenizer struct {
	Tokens Tokens
}

type Tokens struct {
	Values []Token `json:"values"`
}

type Token struct {
	Keyword string `json:"keyword"`
	Tag     string `json:"tag"`
	Tf      int    `json:"tf"`
}

func (t *Tokenizer) Tokenize(text string, lang string) (result Tokens, err error) {
	client := &http.Client{
		Timeout: time.Millisecond * 60000,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	if lang != "ja" && lang != "ko" && lang != "cn" && lang != "en" {
		return result, errors.New("lang [" + lang + "] is not supported.")
	}

	postbody := strings.NewReader(fmt.Sprintf("text=%s", text))
	req, err := http.NewRequest("POST", "http://localhost:9100/"+lang, postbody)

	if err != nil {
		return result, err
	}

	res, err := client.Do(req)

	if err != nil {
		return result, err
	}

	rawResultBytes, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return result, err
	}

	t.Tokens = Tokens{}
	json.Unmarshal(rawResultBytes, &t.Tokens)

	//ioutil.WriteFile("./result.txt", rawResultBytes, 0777)

	return t.Tokens, err
}

func Clearning(tokens *Tokens) (err error) {
	pt, err := regexp.Compile(`(@|¥|#|\.|\:|\*|/|-|\=|\+|\~|\^|\]|\[)+`)

	if err != nil {
		return
	}

	for i, token := range tokens.Values {
		repd := pt.ReplaceAll([]byte(token.Keyword), []byte(""))
		//log.Println(string(repd))
		token.Keyword = *(*string)(unsafe.Pointer(&repd))
		tokens.Values[i] = token
	}

	return
}

func IsNoun(lang string, token Token) bool {
	switch lang {
	case "ja":
		return _isJapaneseNoun(token)
	case "ko":
		return _isKoreanNoun(token)
	case "cn":
		return _isChineseNoun(token)
	case "en":
		return _isEnglishNoun(token)
	default:
		return false
	}
	return false
}

func _isJapaneseNoun(token Token) bool {
	tag := token.Tag
	parts := strings.Split(tag, ",")

	switch parts[0] {
	case "名詞":
		return true
	default:
		return false
	}
	return false
}

func _isKoreanNoun(token Token) bool {
	tag := token.Tag
	if strings.HasPrefix(tag, "N") == true {
		return true
	} else {
		return false
	}
}

func _isChineseNoun(token Token) bool {
	tag := token.Tag
	if strings.HasPrefix(tag, "N") == true {
		return true
	} else {
		return false
	}
}

func _isEnglishNoun(token Token) bool {
	tag := token.Tag
	if strings.HasPrefix(tag, "N") == true {
		return true
	} else {
		return false
	}
}
