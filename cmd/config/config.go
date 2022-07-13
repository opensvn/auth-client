package config

import "github.com/opensvn/auth-client"

type MqttConfig struct {
	ServerAddr        string `yaml:"server_addr"`         // MQTT server URL
	ClientID          string `yaml:"client_id"`           // client id to use when connecting to server
	ClientName        string `yaml:"client_name"`         // client name
	Topic             string `yaml:"topic"`               // topic on which to publish messaged
	Qos               byte   `yaml:"qos"`                 // qos to use when publishing
	Keepalive         uint16 `yaml:"keepalive"`           // seconds between keepalive packets
	ConnectRetryDelay uint16 `yaml:"connect_retry_delay"` // period between connection attempts
	WriteToStdOut     bool   `yaml:"write_to_stdout"`     // If true received messages will be written to stdout
	WriteToDisk       bool   `yaml:"write_to_disk"`       // if true received messages will be written to below file
	OutputFileName    string `yaml:"output_filename"`     // filename to save messages to
	Debug             bool   `yaml:"debug"`               // autopaho and paho debug output requested
}

type AddrConfig struct {
	Ra       string `yaml:"ra"`
	Platform string `yaml:"platform"`
}

// Config holds the configuration
type Config struct {
	Mqtt MqttConfig        `yaml:"mqtt"`
	User client.UserConfig `yaml:"user"`
	Addr AddrConfig        `yaml:"addr"`
}
