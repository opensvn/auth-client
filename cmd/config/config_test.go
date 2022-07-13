package config

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestParseConfig(t *testing.T) {
	buff, err := ioutil.ReadFile("config.yml")
	assert.Nil(t, err)

	conf := &Config{}
	err = yaml.Unmarshal(buff, conf)
	assert.Nil(t, err)

	assert.NotEmpty(t, conf.Mqtt.ServerAddr)
	assert.NotEmpty(t, conf.User.Uid)
	assert.NotEmpty(t, conf.Addr.Ra)
	assert.NotEmpty(t, conf.Addr.Platform)
}
