package config

import (
	"os"
)

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
	c.twitterApiOauthToken = os.Getenv("WUHAN_TWITTER_API_OAUTH_TOKEN")
	c.twitterApiOauthTokenSecret = os.Getenv("WUHAN_TWITTER_API_OAUTH_TOKEN_SECRET")

}

func (c *Config) MysqlHost() string {
	return c.mysqlHost
}

func (c *Config) MysqlPort() string {
	return c.mysqlPort
}

func (c *Config) MysqlUser() string {
	return c.mysqlUser
}

func (c *Config) MysqlPassword() string {
	return c.mysqlPasswd
}

func (c *Config) MongoHost() string {
	return c.mongoHost
}

func (c *Config) MongoPort() string {
	return c.mongoPort
}

func (c *Config) TwiiterApiConsumerKey() string {
	return c.twitterApiConsumerKey
}

func (c *Config) TwitterApiConsumerSecretKey() string {
	return c.twitterApiConsumerSecretKey
}
