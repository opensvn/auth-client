package main

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

type UserConfig struct {
	Uid                    string `yaml:"uid"`
	Hid                    byte   `yaml:"hid"`
	EncryptPrivateKey      string `yaml:"encrypt_private_key"`
	SignPrivateKey         string `yaml:"sign_private_key"`
	EncryptMasterPublicKey string `yaml:"encrypt_master_public_key"`
	SignMasterPublicKey    string `yaml:"sign_master_public_key"`
}

type RaConfig struct {
	Addr string `yaml:"addr"`
}

// Config holds the configuration
type Config struct {
	Mqtt MqttConfig `yaml:"mqtt"`
	User UserConfig `yaml:"user"`
	Ra   RaConfig   `yaml:"ra"`
}
