package tokenizer

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenizeJapanese(t *testing.T) {
	text := "今日はいい天気ですね"
	var tknizer Tokenizer = Tokenizer{}
	tokens, err := tknizer.Tokenize(text, "ja")

	if err != nil {
		t.Fatal(err.Error())
	}

	for _, v := range tokens.Values {
		t.Log(fmt.Sprintf("%s/%s/%d", v.Keyword, v.Tag, v.Tf))
	}

	t.Fatal("ended")
}

func TestTokenizeKorean(t *testing.T) {
	text := "안녕하세요! 좋은 날씨입니다.오늘은 점심은 뭘 먹을까요?"
	var tknizer Tokenizer = Tokenizer{}
	tokens, err := tknizer.Tokenize(text, "ko")

	if err != nil {
		t.Fatal(err.Error())
	}

	for _, v := range tokens.Values {
		t.Log(fmt.Sprintf("%s/%s/%d", v.Keyword, v.Tag, v.Tf))
	}

	t.Fatal("ended")
}

func TestTokenizeChinese(t *testing.T) {
	text := "湖北记者陈卓摄于武汉新华医院现场。不能告别，遗体直接火化，喊妈妈的声音带着撕心裂肺的哭泣，惨绝人寰！"
	// at the Huoshenshan Hospital, which was finished in 10 days in #Wuhan. On February 2, the makeshift hospital beg"
	var tknizer Tokenizer = Tokenizer{}
	tokens, err := tknizer.Tokenize(text, "cn")

	if err != nil {
		t.Fatal(err.Error())
	}

	for _, v := range tokens.Values {
		t.Log(fmt.Sprintf("%s/%s/%d", v.Keyword, v.Tag, v.Tf))
	}

	t.Fatal("ended")
}

func TestTokenizeEnglish(t *testing.T) {
	text := "An inside look at the Huoshenshan Hospital, which was finished in 10 days in #Wuhan. On February 2, the makeshift hospital beg"
	// at the Huoshenshan Hospital, which was finished in 10 days in #Wuhan. On February 2, the makeshift hospital beg"
	var tknizer Tokenizer = Tokenizer{}
	tokens, err := tknizer.Tokenize(text, "en")

	if err != nil {
		t.Fatal(err.Error())
	}

	for _, v := range tokens.Values {
		t.Log(fmt.Sprintf("%s/%s/%d", v.Keyword, v.Tag, v.Tf))
	}

	t.Fatal("ended")
}

func TestFunctions(t *testing.T) {
	str := "sssss@@¥#.:*/-=+~^]["
	var tk Token = Token{
		Keyword: str,
		Tf:      1,
		Tag:     "tag",
	}

	var tks Tokens = Tokens{
		Values: []Token{
			tk,
		},
	}

	err := Clearning(&tks)
	if err != nil {
		t.Fatal(err.Error())
	}

	assert.Equal(t, tks.Values[0].Keyword, "sssss")
}
