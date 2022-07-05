package main

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestParseConfig(t *testing.T) {
	buff, err := ioutil.ReadFile("config.yml")
	assert.Equal(t, err, nil)

	conf := &Config{}
	err = yaml.Unmarshal(buff, conf)
	assert.Equal(t, err, nil)

	assert.NotEmpty(t, conf.Mqtt.ServerAddr)
	assert.NotEmpty(t, conf.User.Uid)
}
