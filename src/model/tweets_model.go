package model

import (
	"context"
	"fmt"
	"log"
	"ml"
	"time"
	"tokenizer"
	"utils"

	. "config"
	. "sync"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var GlobalConfig Config = Config{}

type TweetModel struct {
	Id         int64
	Username   string
	Screenname string
	Lang       string
	Tweet      string
	CreatedAt  string
}

type TweetTopicModel struct {
	Lang           string
	Topic          string
	Score          float64
	CalculatedDate string
	Id             int64
}

type TweetPickupUsersModel struct {
	Id             int64
	Lang           string
	Screenname     string
	Score          int
	CalculatedDate string
}

type TokenizedTweetsModel struct {
	client   *mongo.Client
	ctx      context.Context
	database *mongo.Database
}

type TokenizedTweets struct {
	TweetId   int64
	Lang      string
	Tokens    []tokenizer.Token
	CreatedAt time.Time
}

type MongoTokenizedTweets struct {
	Id      primitive.ObjectID `json:"_id" bson:"_id"`
	Lang    string             `json:"lang" bson:"lang"`
	TweetId int64              `json:"tweetid" bson:"tweetid"`
	Tokens  []MongoToken       `json:"tokens" bson:"tokens"`
}

type MongoToken struct {
	Keyword string `json:"keyword" bson:"keyword"`
	Tf      int    `json:"tf" bson:"tf"`
	Tag     string `json:"tag" bson:"tag"`
}

type CalculatedTfIdfTokenizedTweets struct {
	TweetId int64
	Lang    string
	TfIdfs  []CalculatedTfIdf
}

type CalculatedTfIdf struct {
	//Word  string
	Hash  uint32
	TfIdf float64
}

type MongoCalculatedTfIdfTokenizedTweets struct {
	TweetId int64                  `json:"tweetid" bson:"tweetid"`
	Lang    string                 `json:"lang" bson:"lang"`
	TfIdfs  []MongoCalculatedTfIdf `json:"tfidfs" bson:"tfidfs"`
}

type MongoCalculatedTfIdf struct {
	Hash  uint32  `json:"hash" bson:"hash"`
	TfIdf float64 `json:"tfidf" bson:"tfidf"`
}

type MongoHashWordMap struct {
	Lang   string                 `json:"lang" bson:"lang"`
	Values MongoHashWordMapValues `json:"hashmap" bson:"hashmap"`
}

type MongoHashWordMapValues struct {
	Values map[string]string `json:"hashmap" bson:"hashmap"`
}

var Db *sqlx.DB

func init() {
	GlobalConfig.Init()
	var err error
	Db, err = sqlx.Connect("mysql", fmt.Sprintf("%s:%s@(%s:%s)/wuhan", GlobalConfig.MysqlUser(), GlobalConfig.MysqlPassword(), GlobalConfig.MysqlHost(), GlobalConfig.MysqlPort()))

	if err != nil {
		log.Fatal(err.Error())
	}
}

func (twMs *TweetModel) clear() {
	twMs.Id = 0
	twMs.Username = ""
	twMs.Lang = ""
	twMs.Tweet = ""
	twMs.CreatedAt = ""
}

func (twTMs *TweetTopicModel) clear() {
	twTMs.Lang = ""
	twTMs.Topic = ""
	twTMs.Score = 0.0
}

func (twPuMs *TweetPickupUsersModel) clear() {
	twPuMs.Lang = ""
	twPuMs.Screenname = ""
	twPuMs.Score = 0.0
}

func (twMs *TweetModel) FindSinceId(lang string) (sinceId int64, err error) {
	stmt, err := Db.Prepare("select ifnull(max(id),0) as since_id from tweets where lang = ? ")
	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	row := stmt.QueryRow(lang)
	err = row.Scan(&sinceId)

	if err != nil {
		return 0, err
	}

	return
}

func (twMs *TweetModel) FindByLang(lang string) (result []TweetModel, err error) {
	stmt, err := Db.Prepare("select id,lang,username,screen_name,tweet,created_at from tweets where lang = ?")
	if err != nil {
		return result, err
	}

	defer stmt.Close()

	rows, err := stmt.Query(lang)

	if err != nil {
		return result, err
	}

	defer rows.Close()

	for rows.Next() {
		var tmpTwMs TweetModel = TweetModel{}
		err := rows.Scan(&tmpTwMs.Id, &tmpTwMs.Lang, &tmpTwMs.Username, &tmpTwMs.Screenname, &tmpTwMs.Tweet, &tmpTwMs.CreatedAt)

		if err != nil {
			return result, err
		}

		result = append(result, tmpTwMs)
	}

	return
}

func (twMs *TweetModel) Insert() {
	defer twMs.clear()

	stmt, err := Db.Prepare("select count(*) as ct from tweets where id = ?")

	if err != nil {
		panic(err.Error())
	}

	defer stmt.Close()

	var ct int = 0
	row := stmt.QueryRow(twMs.Id)
	row.Scan(&ct)

	//log.Printf("ct:%d", ct)
	if ct > 0 {
		return
	}

	Db.MustExec(
		`insert into tweets(id,username,screen_name,lang,tweet,created_at) values(?,?,?,?,?,?)`,
		twMs.Id, twMs.Username, twMs.Screenname, twMs.Lang, twMs.Tweet, twMs.CreatedAt,
	)
}

func (twTMs *TweetTopicModel) Begin() {
	Db.MustBegin()
}

func (twTMs *TweetTopicModel) Commit() {
	Db.MustBegin().Commit()
}

func (twTMs *TweetTopicModel) FindByLang(lang string, targetDate time.Time) (result []TweetTopicModel, err error) {
	stmt, err := Db.Prepare("select id,topic,score,lang,calculated_date from tweet_topics where lang = ? and calculated_date = ? order by score desc")

	if err != nil {
		return result, err
	}

	rows, err := stmt.Query(lang, utils.FormatDate(targetDate))

	if err != nil {
		return result, err
	}

	for rows.Next() {
		var rec TweetTopicModel
		rows.Scan(&rec.Id, &rec.Topic, &rec.Score, &rec.Lang, &rec.CalculatedDate)
		result = append(result, rec)
	}

	return
}

func (twTMs *TweetTopicModel) InsertTopic() {
	defer twTMs.clear()

	Db.MustExec(
		`insert into tweet_topics(lang,topic,score,calculated_date) values(?,?,?,Now())`,
		twTMs.Lang, twTMs.Topic, twTMs.Score,
	)
}

func (twTMs *TweetTopicModel) DeleteTopic(lang string) (err error) {

	_, err = Db.Exec("delete from tweet_topics where lang = ? and calculated_date = ? ", lang, utils.FormatTime(time.Now()))

	if err != nil {
		panic(err.Error())
	}

	return
}

func (twPuMs *TweetPickupUsersModel) InsertPickedUpUser() {
	defer twPuMs.clear()

	Db.MustExec(
		`insert into tweet_pickup_users(screen_name,lang,score,calculated_date) values(?,?,?,Now())`,
		twPuMs.Screenname, twPuMs.Lang, twPuMs.Score,
	)
}

func (twPuMs *TweetPickupUsersModel) DeletePickupUsers(lang string) (err error) {

	_, err = Db.Exec("delete from tweet_pickup_users where lang = ? and calculated_date = ?", lang, utils.FormatTime(time.Now()))

	if err != nil {
		panic(err.Error())
	}

	return
}

func (tpuM *TweetPickupUsersModel) FindByLang(lang string, targetDate time.Time) (result []TweetPickupUsersModel, err error) {
	stmt, err := Db.Prepare("select id,screen_name,score,lang,calculated_date from tweet_pickup_users where lang = ? and calculated_date = ? order by score desc limit 30")

	if err != nil {
		return result, err
	}

	rows, err := stmt.Query(lang, utils.FormatDate(targetDate))

	if err != nil {
		return result, err
	}

	for rows.Next() {
		var rec TweetPickupUsersModel
		rows.Scan(&rec.Id, &rec.Screenname, &rec.Score, &rec.Lang, &rec.CalculatedDate)
		result = append(result, rec)
	}

	return
}

/**
var Mongo *mongo.Client
var Ctx context.Context
**/
func (ttM *TokenizedTweetsModel) Init() {
	Mongo, err := mongo.NewClient(options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", GlobalConfig.MongoHost(), GlobalConfig.MongoPort())))

	if err != nil {
		log.Fatal(err.Error())
	}

	ctx, _ := context.WithTimeout(context.Background(), 300*time.Second)
	err = Mongo.Connect(ctx)

	if err != nil {
		log.Fatal(err.Error())
	}

	ttM.client = Mongo
	ttM.ctx = ctx
	ttM.database = Mongo.Database("wuhan")
}

func (ttM *TokenizedTweetsModel) GetDatabase() *mongo.Database {
	return ttM.database
}

func (ttM *TokenizedTweetsModel) GetClient() mongo.Client {
	return *ttM.client
}

func (ttM *TokenizedTweetsModel) GetCtx() context.Context {
	return ttM.ctx
}

func (ttM *TokenizedTweetsModel) HasTokenizedTweet(tweetId int64) (result bool, err error) {
	collection := ttM.GetDatabase().Collection("tweets")

	ct, err := collection.CountDocuments(ttM.GetCtx(), bson.M{"tweetid": tweetId})

	if err != nil {
		return result, err
	}

	if ct > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func (ttM *TokenizedTweetsModel) InsertTokenizedTweets(tweetId int64, lang string, createdAt string, tokens *tokenizer.Tokens, wg *WaitGroup) {
	defer wg.Done()

	collection := ttM.GetDatabase().Collection("tweets")

	var tokenizedTweets TokenizedTweets = TokenizedTweets{
		TweetId:   tweetId,
		Lang:      lang,
		Tokens:    tokens.Values,
		CreatedAt: utils.ParseStringDate(createdAt),
	}

	var upsert bool = true
	var opt options.UpdateOptions = options.UpdateOptions{}
	opt.Upsert = &upsert

	//for dosent insert same tweet
	_, err := collection.UpdateOne(ttM.ctx, bson.D{{"tweetid", tweetId}}, bson.D{{"$set", tokenizedTweets}}, &opt)
	//_, err = collection.InsertOne(ttM.ctx, &tokenizedTweets)

	if err != nil {
		panic(err.Error())
	}
}

func (ttM *TokenizedTweetsModel) InsertTokenizedTweetsold(tokenizedTweets *[]interface{}) (err error) {
	collection := ttM.GetDatabase().Collection("tweets")

	_, err = collection.InsertMany(ttM.ctx, *tokenizedTweets)

	if err != nil {
		return err
	}

	return err
}

func (ttM *TokenizedTweetsModel) InsertTfIdfTweets(tweetId int64, lang string, tfidfs ml.FloatVector, hashWordMap *ml.HashWordMap, wg *WaitGroup) {
	defer wg.Done()

	collection := ttM.GetDatabase().Collection("tweets_tfidf")
	var tfIdfs []CalculatedTfIdf = make([]CalculatedTfIdf, 0, 100)

	for hashNum, tfIdf := range tfidfs.Values {
		if tfIdf > 0 {
			tfIdfs = append(tfIdfs, CalculatedTfIdf{
				//Word:  hashWordMap.HashMap[hashNum],
				Hash:  hashNum,
				TfIdf: tfIdf,
			})
		}
	}

	var doc = CalculatedTfIdfTokenizedTweets{
		TweetId: tweetId,
		Lang:    lang,
		TfIdfs:  tfIdfs,
	}

	var upsert bool = true
	var opt options.UpdateOptions = options.UpdateOptions{}
	opt.Upsert = &upsert

	_, err := collection.UpdateOne(ttM.ctx, bson.D{{"tweetid", tweetId}}, bson.D{{"$set", doc}}, &opt)
	//_, err = collection.InsertOne(ttM.ctx, &doc)

	if err != nil {
		panic(err.Error())
	}
}

func (ttM *TokenizedTweetsModel) InsertHashWordMap(lang string, hashWordMap *ml.HashWordMap) (err error) {
	collection := ttM.GetDatabase().Collection("tweets_hashwordmap")

	_, err = collection.DeleteOne(ttM.ctx, bson.D{{"lang", lang}})

	if err != nil {
		return err
	}

	_, err = collection.InsertOne(ttM.ctx, bson.D{{"lang", lang}, {"hashmap", &hashWordMap}})

	if err != nil {
		return err
	}

	return err
}

func (ttM *TokenizedTweetsModel) FindHashWordMapByLang(lang string) (cur *mongo.Cursor, err error) {
	collection := ttM.GetDatabase().Collection("tweets_hashwordmap")
	cur, err = collection.Find(ttM.ctx, bson.M{"lang": lang})

	return
}

func (ttM *TokenizedTweetsModel) FindByLang(lang string, rangeStart interface{}) (cur *mongo.Cursor, err error) {
	collection := ttM.GetDatabase().Collection("tweets")

	var cond bson.M
	if rangeStart != nil {
		cond = bson.M{"lang": lang, "createdat": bson.M{"$gt": rangeStart}}
	} else {
		cond = bson.M{"lang": lang}
	}
	cur, err = collection.Find(ttM.ctx, cond)

	return
}

func (ttM *TokenizedTweetsModel) FindByTweetId(tweetId int64) (cur *mongo.Cursor, err error) {
	collection := ttM.GetDatabase().Collection("tweets")
	cur, err = collection.Find(ttM.ctx, bson.M{"tweetid": tweetId})

	return
}

func (ttM *TokenizedTweetsModel) FindTfIdfsByLang(lang string) (cur *mongo.Cursor, err error) {
	collection := *ttM.GetDatabase().Collection("tweets_tfidf")
	cur, err = collection.Find(ttM.ctx, bson.M{"lang": lang})

	return
}

func (ttM *TokenizedTweetsModel) FindTfIdfByTweetId(tweetId int64) (cur *mongo.Cursor, err error) {
	collection := *ttM.GetDatabase().Collection("tweets_tfidf")
	cur, err = collection.Find(ttM.ctx, bson.M{"tweetid": tweetId})

	return
}
