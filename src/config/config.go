package config

import "os"

type Config struct {
	mysqlHost                   string
	mysqlPort                   string
	mysqlUser                   string
	mysqlPasswd                 string
	mongoHost                   string
	mongoPort                   string
	twitterApiConsumerKey       string
	twitterApiConsumerSecretKey string
	twitterApiOauthToken        string
	twitterApiOauthTokenSecret  string
}

func (c *Config) Init() {
	c.mysqlHost = os.Getenv("WUHAN_MYSQL_HOST")
	c.mysqlPort = os.Getenv("WUHAN_MYSQL_PORT")
	c.mysqlUser = os.Getenv("WUHAN_MYSQL_USER")
	c.mysqlPasswd = os.Getenv("WUHAN_MYSQL_PASS")
	c.mongoHost = os.Getenv("WUHAN_MONGO_HOST")
	c.mongoPort = os.Getenv("WUHAN_MONGO_PORT")
	c.twitterApiConsumerKey = os.Getenv("WUHAN_TWITTER_API_CONSUMER_KEY")
	c.twitterApiConsumerSecretKey = os.Getenv("WUHAN_TWITTER_API_CONSUMER_SECRET_KEY")
	
}

func (c *Config) MysqlHost() string {

}