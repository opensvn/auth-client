package main

import (
	"encoding/hex"
	"log"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/opensvn/auth-client"
	"github.com/opensvn/auth-client/cmd/config"
	"github.com/opensvn/auth-client/log"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

var globalUser *client.User

func main() {
	buf, err := os.ReadFile("config/config.yml")
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

	logging.InitLogger(conf.Log.Level, conf.Log.Path)

	globalUser = client.NewUser(&conf.User)
	if globalUser == nil || globalUser.GetEncryptPrivateKey() == nil || globalUser.GetSignPrivateKey() == nil {
		random, err := ApplyKey(conf, globalUser)
		if err != nil {
			logging.Logger.Error("apply key error", zap.Error(err))
			return
		}

		wg := &sync.WaitGroup{}
		wg.Add(1)
		go func() {
			for {
				time.Sleep(time.Second * 3)
				keys, err := queryKey(conf)
				if err != nil {
					logging.Logger.Error("query key", zap.Error(err))
					continue
				}

				if keys.EncryptKey == "" || keys.SignKey == "" {
					continue
				}

				signKeyBuf, err := hex.DecodeString(keys.SignKey)
				if err != nil {
					logging.Logger.Error("decode string", zap.Error(err))
					continue
				}

				encryptKeyBuf, err := hex.DecodeString(keys.EncryptKey)
				if err != nil {
					logging.Logger.Error("decode string", zap.Error(err))
					continue
				}

				encryptKey, err := client.OfbEncrypt(random, encryptKeyBuf)
				if err != nil {
					logging.Logger.Error("encrypt", zap.Error(err))
					continue
				}

				signKey, err := client.OfbEncrypt(random, signKeyBuf)
				if err != nil {
					logging.Logger.Error("encrypt", zap.Error(err))
					continue
				}

				conf.User.EncryptPrivateKey = hex.EncodeToString(encryptKey)
				conf.User.SignPrivateKey = hex.EncodeToString(signKey)
				break
			}
			wg.Done()
		}()
		wg.Wait()

		err = globalUser.SetEncryptPrivateKey(conf.User.EncryptPrivateKey)
		if err != nil {
			logging.Logger.Error("set encrypt private key error", zap.Error(err))
			return
		}

		err = globalUser.SetSignPrivateKey(conf.User.SignPrivateKey)
		if err != nil {
			logging.Logger.Error("set sign private key error", zap.Error(err))
			return
		}

		// save yml file
		buf, err := yaml.Marshal(conf)
		if err != nil {
			logging.Logger.Error("marshal failed", zap.Error(err))
			return
		}

		err = os.WriteFile("config/config.yml", buf, 0644)
		if err != nil {
			logging.Logger.Error("write config file error", zap.Error(err))
			return
		}
	}

	serverUrl, err := url.Parse(conf.Mqtt.ServerAddr)
	if err != nil {
		logging.Logger.Error("url parse", zap.Error(err))
		return
	}

	c := &client.Client{
		Config: &client.ClientConfig{
			Name:              conf.Mqtt.ClientName,
			Topic:             conf.Mqtt.Topic,
			Qos:               conf.Mqtt.Qos,
			Keepalive:         conf.Mqtt.Keepalive,
			ConnectRetryDelay: conf.Mqtt.ConnectRetryDelay,
			Debug:             conf.Mqtt.Debug,
		},
	}
	c.User = globalUser
	c.ServerUrl = serverUrl
	c.AuthHandler = client.NewSm9Auth(c)
	c.SetMsgHandler(HandleMsg)

	// Connect to the broker
	err = c.Connect()
	if err != nil {
		logging.Logger.Error("connect", zap.Error(err))
		return
	}

	quit := make(chan struct{})
	go func() {
		// Publish a message to server periodly
		for {
			select {
			case <-quit:
				return
			default:
				time.Sleep(time.Second * 30)
				_ = c.Publish(conf.Mqtt.Topic, "test")
			}
		}
	}()

	// Messages will be handled through the callback so we really just need to wait until a shutdown is requested
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	signal.Notify(sig, syscall.SIGTERM)

	<-sig
	close(quit)

	// We could cancel the context at this point but will call Disconnect instead (this waits for autopaho to shutdown)
	err = c.Disconnect()
	if err != nil {
		logging.Logger.Error("disconnect", zap.Error(err))
	}
}
