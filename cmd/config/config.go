package config

import "github.com/opensvn/auth-client"

type MqttConfig struct {
	ServerAddr        string `yaml:"server_addr"`         // MQTT server URL
	ClientName        string `yaml:"client_name"`         // client name
	DeviceType        string `yaml:"device_type"`         // device type
	Topic             string `yaml:"topic"`               // topic on which to publish messaged
	Qos               byte   `yaml:"qos"`                 // qos to use when publishing
	Keepalive         uint16 `yaml:"keepalive"`           // seconds between keepalive packets
	ConnectRetryDelay uint16 `yaml:"connect_retry_delay"` // period between connection attempts
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
