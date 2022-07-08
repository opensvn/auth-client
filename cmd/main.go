package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/opensvn/auth-client"
	"github.com/opensvn/auth-client/cmd/config"
	"gopkg.in/yaml.v3"
)

func main() {
	buf, err := ioutil.ReadFile("config/config.yml")
	if err != nil {
		panic(err)
	}

	conf := &config.Config{}
	err = yaml.Unmarshal(buf, conf)
	if err != nil {
		panic(err)
	}

	user := InitUser(conf)

	serverUrl, err := url.Parse(conf.Mqtt.ServerAddr)
	if err != nil {
		log.Printf("%s\n", err)
		return
	}

	c := &client.Client{
		Config: &client.ClientConfig{
			ClientID: conf.Mqtt.ClientID,
			ClientName: conf.Mqtt.ClientName,
			Topic: conf.Mqtt.Topic,
			Qos: conf.Mqtt.Qos,
			Keepalive: conf.Mqtt.Keepalive,
			ConnectRetryDelay: conf.Mqtt.ConnectRetryDelay,
			WriteToStdOut: conf.Mqtt.WriteToStdOut,
			WriteToDisk: conf.Mqtt.WriteToDisk,
			OutputFileName: conf.Mqtt.OutputFileName,
			Debug: conf.Mqtt.Debug,
		},
	}
	c.User = user
	c.ServerUrl = serverUrl
	c.AuthHandler = client.NewSm9Auth(c)

	// Connect to the broker
	err = c.Connect()
	if err != nil {
		log.Printf("%s\n", err)
		return
	}

	// Messages will be handled through the callback so we really just need to wait until a shutdown is requested
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	signal.Notify(sig, syscall.SIGTERM)

	<-sig
	fmt.Println("signal caught - exiting")

	// We could cancel the context at this point but will call Disconnect instead (this waits for autopaho to shutdown)
	err = c.Disconnect()
	if err != nil {
		log.Printf("%s\n", err)
	}

	fmt.Println("shutdown complete")
}

func InitUser(conf *config.Config) *client.User {
	user := &client.User{
		Uid: []byte(conf.User.Uid),
		Hid: conf.User.Hid,
	}

	err := user.SetEncryptPrivateKey(conf.User.EncryptPrivateKey)
	if err != nil {
		panic(err)
	}

	err = user.SetEncryptMasterPublicKey(conf.User.EncryptMasterPublicKey)
	if err != nil {
		panic(err)
	}

	err = user.SetSignPrivateKey(conf.User.SignPrivateKey)
	if err != nil {
		panic(err)
	}

	err = user.SetSignMasterPublicKey(conf.User.SignMasterPublicKey)
	if err != nil {
		panic(err)
	}

	return user
}
