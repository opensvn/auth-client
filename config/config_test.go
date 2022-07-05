package config

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"testing"
)

func TestParseConfig(t *testing.T) {
	buff, err := ioutil.ReadFile("config.yml")
	assert.Equal(t, err, nil)

	conf := &Config{}
	err = yaml.Unmarshal(buff, conf)
	assert.Equal(t, err, nil)
}
