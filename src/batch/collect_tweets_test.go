package batch

import (
	"testing"
)

func TestCollectTweets(t *testing.T) {
	err := Collect_tweet()

	if err != nil {
		t.Fatal(err.Error())
	}

	//	t.Fatal("ended")
}
