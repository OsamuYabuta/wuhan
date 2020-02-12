package ml

import (
	"errors"
	//"log"
	"math/rand"
	"time"
)

type CalculatedTfIdfForTopicModel struct {
	Hash  uint32
	TfIdf float64
}

//simple topic model
//http://www.phontron.com/slides/nlp-programming-ja-07-topic.pdf
func CalculateTopic(calculateTfIdfs *[]CalculatedTfIdfForTopicModel) ([]uint32, map[uint32]float64, error) {
	//sum of tf-idf of all words
	var sumOfTfIdfEachWord map[uint32]float64 = make(map[uint32]float64)
	var sumOfTfIdfInAllWord float64 = 0.0
	//log.Printf("%d", len(*calculateTfIdfs))
	for _, tfidfData := range *calculateTfIdfs {
		hashNum := tfidfData.Hash
		tfIdf := tfidfData.TfIdf + 1.0 //for minus to plus
		sumOfTfIdfEachWord[hashNum] += tfIdf
		sumOfTfIdfInAllWord += tfIdf
	}

	//var probsEachWord map[uint32]float64 = make(map[uint32]float64)
	var hashNumList []uint32 = make([]uint32, 0, 1000000)
	var probsList []float64 = make([]float64, 0, 1000000)
	for _, tfidfData := range *calculateTfIdfs {
		hashNum := tfidfData.Hash
		hashNumList = append(hashNumList, hashNum)
		prob := sumOfTfIdfEachWord[hashNum] / sumOfTfIdfInAllWord
		probsList = append(probsList, prob)
		//log.Printf("%f", probsEachWord[hashNum])
	}

	var idxesIncludeHashNum map[uint32][]int = make(map[uint32][]int)

	for idx, hashNum := range hashNumList {
		idxesIncludeHashNum[hashNum] = append(idxesIncludeHashNum[hashNum], idx)
	}

	var pickedUpList []int = make([]int, 0, 30)
	var pickedUpScoreList []float64 = make([]float64, 0, 30)
	var repeat int = 0
	pickedUpList, pickedUpScoreList, err := _pickUpTopic(probsList, pickedUpList, pickedUpScoreList, repeat, &idxesIncludeHashNum, &hashNumList)

	if err != nil {
		return []uint32{0}, map[uint32]float64{}, err
	}

	var pickedUpHashNumList []uint32 = make([]uint32, 0, 30)
	var pickedUpHashNumScoreMap map[uint32]float64 = make(map[uint32]float64)

	for _idx, idx := range pickedUpList {
		pickedUpHashNumList = append(pickedUpHashNumList, hashNumList[idx])
		pickedUpHashNumScoreMap[hashNumList[idx]] = pickedUpScoreList[_idx]
	}

	return pickedUpHashNumList, pickedUpHashNumScoreMap, err
}

func _hasIdx(idx int, pickedUpList []int, idxesIncludeHashNum *map[uint32][]int, hashNumList *[]uint32) bool {
	for _, _idx := range pickedUpList {
		if _idx == idx {
			return true
		}

		hashNum := (*hashNumList)[_idx]
		idxes := (*idxesIncludeHashNum)[hashNum]

		//log.Println(idxes)
		for _, idx2 := range idxes {
			if idx2 == idx {
				return true
			}
		}
	}

	return false
}

func _pickUpTopic(probList []float64, pickedUpList []int, pickedUpScoreList []float64, repeat int, idxesIncludeHashNum *map[uint32][]int, hashNumList *[]uint32) ([]int, []float64, error) {
	var sumOfProps float64 = 0.0

	for idx, prob := range probList {
		//log.Printf("%f", prob)
		if _hasIdx(idx, pickedUpList, idxesIncludeHashNum, hashNumList) {
			continue
		}
		sumOfProps += prob
	}

	rand.Seed(time.Now().UnixNano())
	var remaining = rand.Float64() + 1.0
	//log.Printf("%f", remaining)
	for idx, prob := range probList {
		if _hasIdx(idx, pickedUpList, idxesIncludeHashNum, hashNumList) {
			//fmt.Println("idx:%d is skip", idx)
			continue
		}

		remaining -= prob
		if remaining <= 0 {
			repeat += 1

			if repeat >= 30 {
				return pickedUpList, pickedUpScoreList, nil
			} else {
				pickedUpList = append(pickedUpList, idx)
				pickedUpScoreList = append(pickedUpScoreList, prob)
				return _pickUpTopic(probList, pickedUpList, pickedUpScoreList, repeat, idxesIncludeHashNum, hashNumList)
			}
		}
	}

	if len(pickedUpList) > 0 {
		return pickedUpList, pickedUpScoreList, nil
	} else {
		return pickedUpList, pickedUpScoreList, errors.New("failed to pick")
	}
}
