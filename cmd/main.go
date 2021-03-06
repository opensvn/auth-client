package main

import (
	"encoding/hex"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/opensvn/auth-client"
	"github.com/opensvn/auth-client/cmd/config"
	"gopkg.in/yaml.v3"
)

func main() {
	buf, err := ioutil.ReadFile("config/config.yml")
	if err != nil {
		log.Printf("read config file error: %v\n", err)
		return
	}

	conf := &config.Config{}
	err = yaml.Unmarshal(buf, conf)
	if err != nil {
		log.Printf("unmarshal config error: %v\n", err)
		return
	}

	user := client.NewUser(&conf.User)
	if user == nil || user.GetEncryptPrivateKey() == nil || user.GetSignPrivateKey() == nil {
		random, err := ApplyKey(conf, user)
		if err != nil {
			log.Printf("apply key error: %v\n", err)
			return
		}

		wg := &sync.WaitGroup{}
		wg.Add(1)
		go func () {
			for {
				time.Sleep(time.Second * 3)
				keys, err := queryKey(conf)
				if err != nil {
					log.Println(err)
					continue
				}

				if keys.EncryptKey == "" || keys.SignKey == "" {
					continue
				}

				signKeyBuf, err := hex.DecodeString(keys.SignKey)
				if err != nil {
					log.Println(err)
					continue
				}

				encryptKeyBuf, err := hex.DecodeString(keys.EncryptKey)
				if err != nil {
					log.Println(err)
					continue
				}

				encryptKey, err := OfbEncrypt(random, encryptKeyBuf)
				if err != nil {
					log.Println(err)
					continue
				}

				signKey, err := OfbEncrypt(random, signKeyBuf)
				if err != nil {
					continue
				}

				conf.User.EncryptPrivateKey = hex.EncodeToString(encryptKey)
				conf.User.SignPrivateKey = hex.EncodeToString(signKey)
				break
			}
			wg.Done()
		}()
		wg.Wait()

		err = user.SetEncryptPrivateKey(conf.User.EncryptPrivateKey)
		if err != nil {
			log.Printf("set encrypt private key error: %v\n", err)
			return
		}

		err = user.SetSignPrivateKey(conf.User.SignPrivateKey)
		if err != nil {
			log.Printf("set sign private key error: %v\n", err)
			return
		}

		// save yml file
		buf, err := yaml.Marshal(conf)
		if err != nil {
			log.Printf("marshal failed: %v\n", err)
			return
		}

		err = ioutil.WriteFile("config/config.yml", buf, 0644)
		if err != nil {
			log.Printf("write config file error: %v\n", err)
			return
		}
	}

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

	// We could cancel the context at this point but will call Disconnect instead (this waits for autopaho to shutdown)
	err = c.Disconnect()
	if err != nil {
		log.Printf("%s\n", err)
	}
}
