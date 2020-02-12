package ml

import (
	"hash/fnv"
	//"log"
	"math"
	"time"
)

type Matrix struct {
	Rows []Vector
}

type Vector struct {
	WordTf []WordTf
}

type WordTf struct {
	Word string
	Tf   int
	Hash uint32
}

type WorkMatrix1 struct {
	Rows []map[uint32]Tf
}

type Tf struct {
	Tf float64
}

type ResultMatrix struct {
	Rows []FloatVector
}

type FloatVector struct {
	Values map[uint32]float64
}

type HashWordMap struct {
	HashMap map[uint32]string
}

func hash(word string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(word))
	return h.Sum32()
}

func CalculateTfIdf(matrix Matrix) (result ResultMatrix, hashWordMap HashWordMap, err error) {
	//sum of count of tweets
	docCt := len(matrix.Rows)
	//sum of count of words each word in all tweets
	var docsCountIncludeWord map[uint32]int = make(map[uint32]int, 100000)
	//initializer work matrix
	var wm WorkMatrix1 = WorkMatrix1{
		Rows: make([]map[uint32]Tf, docCt),
	}

	hashWordMap.HashMap = make(map[uint32]string, 100000)
	ct := 0
	for i, row := range matrix.Rows {
		//sum of count of words in a tweet
		var allWordTf int = 0
		var wordsTfInThisTweets map[uint32]int = make(map[uint32]int, 100)
		for i2, vector := range row.WordTf {
			hashNum := hash(vector.Word)
			hashWordMap.HashMap[hashNum] = vector.Word
			vector.Hash = hashNum
			matrix.Rows[i].WordTf[i2] = vector
			wordsTfInThisTweets[hashNum] += vector.Tf
			allWordTf += vector.Tf
		}

		wm.Rows[i] = make(map[uint32]Tf, len(wordsTfInThisTweets))
		for hashNum, wordTf := range wordsTfInThisTweets {
			tf := float64(wordTf) / float64(allWordTf)
			wm.Rows[i][hashNum] = Tf{
				Tf: tf,
			}
			docsCountIncludeWord[hashNum] = 1
		}

		if ct%1000 == 0 {
			time.Sleep(time.Millisecond * 200)
		}
		ct++
	}

	//calculate idf
	var resultMatrix ResultMatrix = ResultMatrix{
		Rows: make([]FloatVector, docCt),
	}
	for docIdx, row := range wm.Rows {
		resultMatrix.Rows[docIdx] = FloatVector{
			Values: make(map[uint32]float64, len(docsCountIncludeWord)),
		}

		for hashNum, docCtIncludeWord := range docsCountIncludeWord {
			tf := row[hashNum]
			idf := float64(docCt) / float64(docCtIncludeWord)

			//log.Printf("%f/%f", tf	, idf)
			if tf.Tf == 0 {
				resultMatrix.Rows[docIdx].Values[hashNum] = 0.0
			} else {
				resultMatrix.Rows[docIdx].Values[hashNum] = math.Log(tf.Tf * idf)
			}
		}
	}

	//normalize tf-idf by norm2
	normalize(&resultMatrix)

	result = resultMatrix

	return
}

func normalize(matrix *ResultMatrix) {
	for rowIdx, row := range matrix.Rows {
		//l2 norm
		var powSum float64 = 0.0
		for _, tfIdf := range row.Values {
			powSum += math.Pow(tfIdf, 2)
		}
		norm := math.Sqrt(powSum)
		for hashNum, tfIdf := range row.Values {
			matrix.Rows[rowIdx].Values[hashNum] = tfIdf / norm
		}
	}
}
