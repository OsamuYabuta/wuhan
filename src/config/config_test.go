package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	var config Config = Config{}
	config.Init()

	assert.Equal(t , "localhost" , config.MysqlHost())
	assert.Equal(t , "3306" , config.MysqlPort())
	assert.Equal(t , "root" , config.MysqlUser())
	assert.Equal(t , "osamu2009" , config.MysqlPassword())
}