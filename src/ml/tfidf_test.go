package ml

import (
	"model"
	"testing"
)

func TestTfIdf(t *testing.T) {
	var mt Matrix = Matrix{
		Rows: []Vector{
			Vector{
				WordTf: []WordTf{
					WordTf{
						Word: "ワード1",
						Tf:   1,
						Hash: 0,
					},
					WordTf{
						Word: "ワード2",
						Tf:   1,
						Hash: 0,
					},
					WordTf{
						Word: "ワード3",
						Tf:   2,
						Hash: 0,
					},
					WordTf{
						Word: "ワード4",
						Tf:   10,
						Hash: 0,
					},
					WordTf{
						Word: "ワード5",
						Tf:   15,
						Hash: 0,
					},
				},
			},
			Vector{
				WordTf: []WordTf{
					WordTf{
						Word: "ワード1",
						Tf:   10,
						Hash: 0,
					},
					WordTf{
						Word: "ワード2",
						Tf:   1,
						Hash: 0,
					},
					WordTf{
						Word: "ワード3",
						Tf:   100,
						Hash: 0,
					},
					WordTf{
						Word: "ワード4",
						Tf:   2,
						Hash: 0,
					},
					WordTf{
						Word: "ワード5",
						Tf:   2,
						Hash: 0,
					},
					WordTf{
						Word: "ワード6",
						Tf:   200,
						Hash: 0,
					},
				},
			},
		},
	}
	res, hashWordMap, err := CalculateTfIdf(mt)

	if err != nil {
		t.Fatal(err.Error())
	}

	t.Log(hashWordMap)
	t.Fatal(res)
}

func TestTweetTfIdf(t *testing.T) {
	var ttM model.TokenizedTweetsModel = model.TokenizedTweetsModel{}
	ttM.Init()
	cursor, err := ttM.FindByLang("ja")

	if err != nil {
		t.Fatal(err.Error())
	}

	defer cursor.Close(ttM.GetCtx())
	var results []model.MongoTokenizedTweets
	ct := 0
	for cursor.Next(ttM.GetCtx()) {
		var result model.MongoTokenizedTweets
		err := cursor.Decode(&result)

		if err != nil {
			t.Fatal(err.Error())
		}

		results = append(results, result)

		ct += 1
		if ct > 50 {
			break
		}
	}

	var matrix Matrix = Matrix{
		Rows: make([]Vector, len(results)),
	}
	for i, result := range results {
		matrix.Rows[i].WordTf = make([]WordTf, len(result.Tokens))
		//t.Fatal(result["tokens"])

		for i2, token := range result.Tokens {
			matrix.Rows[i].WordTf[i2] = WordTf{
				Word: token.Keyword,
				Tf:   token.Tf,
			}
		}
	}

	//t.Log("checkpoint1")
	res, hashWordMap, err := CalculateTfIdf(matrix)

	t.Log(res)
	if err != nil {
		t.Fatal(err.Error())
	}

	t.Log(hashWordMap)
	t.Log(res)
	t.Fatal("ended")
}
