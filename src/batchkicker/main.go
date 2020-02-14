package main

import (
	"os"
	//"testing"
	"log"

	. "batch"
)

func main() {
	args := os.Args

	var err error

	switch args[1] {
	case "tokenize":
		err = Tokenize_tweets()
	case "collect":
		err = Collect_tweet()
	case "calc_topic":
		err = Calculate_topic()
	case "calc_pickedup_user":
		err = Calculate_User()
	case "calc_tfidf":
		err = Calculate_tfidf()
	}

	if err != nil {
		log.Println(err.Error())
	}

}
