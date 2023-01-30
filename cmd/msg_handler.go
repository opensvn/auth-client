package main

import (
	"encoding/hex"

	"github.com/eclipse/paho.golang/paho"
	"github.com/opensvn/auth-client"
	"github.com/opensvn/auth-client/log"
	"go.uber.org/zap"
)

func HandleMsg(msg *paho.Publish) {
	payload, err := hex.DecodeString(string(msg.Payload))
	if err != nil {
		logging.Logger.Error("Decode string", zap.Error(err))
		return
	}

	text, err := client.OfbEncrypt(globalUser.SessionKey, payload)
	if err != nil {
		logging.Logger.Error("Decrypt payload", zap.Error(err))
		return
	}

	logging.Logger.Info("Receive ", zap.String("topic", msg.Topic), zap.String("msg", string(text)))
}
