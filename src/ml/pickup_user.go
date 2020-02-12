package ml

import (
	"log"
	"sort"
)

type TokenizedTweetsForPickupUser struct {
	Keyword string
	Tag     string
	Tf      int
}

type KeyValue struct {
	Key   string
	Value int
}

type KeyValues []KeyValue

func (kvs KeyValues) Len() int {
	return len(kvs)
}

func (kvs KeyValues) Swap(i, j int) {
	kvs[i], kvs[j] = kvs[j], kvs[i]
}

func (kvs KeyValues) Less(i, j int) bool {
	return kvs[i].Value > kvs[j].Value
}

func CalculatePickupUser(userNames *map[string]bool, tokens *[]TokenizedTweetsForPickupUser) (result KeyValues) {
	log.Println("start pickup user")

	var userNameTfs map[string]int = make(map[string]int)
	var tfs []int = make([]int, 0, 100000)
	var totalTf int = 0
	//log.Println(*userNames)
	for _, token := range *tokens {
		_, ok := (*userNames)[token.Keyword]
		//log.Println(ok)
		if ok == true {
			userNameTfs[token.Keyword] += token.Tf
			totalTf += token.Tf
		}
	}

	for _, tf := range userNameTfs {
		tfs = append(tfs, tf)
	}

	var percentile4 float32 = 0.0
	sort.Ints(tfs)

	var totalCtF float32 = float32(len(tfs))
	for idx, tf := range tfs {
		if float32(idx)/totalCtF >= 0.75 {
			percentile4 = float32(tf)
			break
		}
	}

	//log.Printf("%d", percentile4)

	var keyValues = KeyValues{}
	for keyword, tf := range userNameTfs {
		if float32(tf) > (percentile4 * 2.5) {
			continue
		}
		keyValue := KeyValue{
			Key:   keyword,
			Value: tf,
		}
		keyValues = append(keyValues, keyValue)
	}

	sort.Sort(keyValues)

	result = keyValues

	return

	//log.Fatal(keyValues[:10])
}
