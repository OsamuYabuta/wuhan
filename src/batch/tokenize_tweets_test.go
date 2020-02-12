package batch

import (
	"testing"
)

func TestTokenizeTweets(t *testing.T) {
	err := Tokenize_tweets()

	if err != nil {
		t.Fatal(err.Error())
	}
}
