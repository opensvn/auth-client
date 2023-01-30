package client

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/url"
	"time"

	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"github.com/opensvn/auth-client/log"
	"go.uber.org/zap"
)

type ClientConfig struct {
	Name              string
	Topic             string
	Qos               byte
	Keepalive         uint16
	ConnectRetryDelay uint16
	Debug             bool
}

type Client struct {
	ServerUrl   *url.URL
	AuthHandler *Sm9Auth
	User        *User
	Cm          *autopaho.ConnectionManager
	Cancel      context.CancelFunc
	Config      *ClientConfig
	MsgHandler  paho.MessageHandler
}

func (c *Client) SetMsgHandler(handler paho.MessageHandler) {
	c.MsgHandler = handler
}

func (c *Client) Connect() error {
	cliCfg := autopaho.ClientConfig{
		BrokerUrls:        []*url.URL{c.ServerUrl},
		KeepAlive:         c.Config.Keepalive,
		ConnectRetryDelay: time.Duration(c.Config.ConnectRetryDelay) * time.Millisecond,
		OnConnectionUp: func(cm *autopaho.ConnectionManager, connAck *paho.Connack) {
			if _, err := cm.Subscribe(context.Background(), &paho.Subscribe{
				Subscriptions: map[string]paho.SubscribeOptions{
					c.Config.Topic: {QoS: c.Config.Qos},
				},
			}); err != nil {
				logging.Logger.Error("failed to subscribe", zap.Error(err))
				return
			}
		},
		OnConnectError: func(err error) {
			logging.Logger.Error("connect", zap.Error(err))
		},
		ClientConfig: paho.ClientConfig{
			ClientID: string(c.User.Uid),
			Router: paho.NewSingleHandlerRouter(func(m *paho.Publish) {
				if c.MsgHandler != nil {
					c.MsgHandler(m)
				}
			}),
			OnClientError: func(err error) {
				logging.Logger.Error("disconnect", zap.Error(err))
			},
			OnServerDisconnect: func(d *paho.Disconnect) {
				if d.Properties != nil {
					logging.Logger.Info("disconnect", zap.String("reason", d.Properties.ReasonString))
				} else {
					logging.Logger.Info("disconnect", zap.Int("code", int(d.ReasonCode)))
				}
			},
			AuthHandler: c.AuthHandler,
		},
	}

	cliCfg.SetConnectPacketConfigurator(func(connect *paho.Connect) *paho.Connect {
		connect.Properties = &paho.ConnectProperties{
			AuthMethod: "sm9",
			AuthData:   []byte(c.AuthHandler.GetRandom1(8)),
			User: []paho.UserProperty{
				{
					Key:   "uid",
					Value: string(c.User.Uid),
				},
				{
					Key:   "hid",
					Value: hex.EncodeToString([]byte{c.User.Hid}),
				},
				{
					Key:   "deviceName",
					Value: c.Config.Name,
				},
			},
		}

		return connect
	})

	if c.Config.Debug {
		cliCfg.Debug = logger{}
		cliCfg.PahoDebug = logger{}
	}

	ctx, cancel := context.WithCancel(context.Background())
	c.Cancel = cancel

	connection, err := autopaho.NewConnection(ctx, cliCfg)
	if err != nil {
		return err
	}

	c.Cm = connection
	return nil
}

func (c *Client) Subscribe(topic string) error {
	subPacket := &paho.Subscribe{
		Subscriptions: map[string]paho.SubscribeOptions{
			topic: {QoS: c.Config.Qos},
		},
	}
	if _, err := c.Cm.Subscribe(context.Background(), subPacket); err != nil {
		logging.Logger.Error("failed to subscribe", zap.Error(err))
		return err
	}

	logging.Logger.Info("subscribe to", zap.String("topic", topic))
	return nil
}

func (c *Client) Publish(topic, payload string) error {
	encPayload, err := OfbEncrypt(c.User.SessionKey, []byte(payload))
	if err != nil {
		logging.Logger.Error("encrypt", zap.Error(err))
		return err
	}

	pubPacket := &paho.Publish{
		Topic:   topic,
		QoS:     byte(0),
		Payload: []byte(hex.EncodeToString(encPayload)),
	}

	if _, err := c.Cm.Publish(context.Background(), pubPacket); err != nil {
		logging.Logger.Error("publish", zap.Error(err))
		return err
	}

	return nil
}

func (c *Client) Disconnect() error {
	defer c.Cancel()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return c.Cm.Disconnect(ctx)
}

type logger struct {
}

func (l logger) Println(v ...interface{}) {
	s := fmt.Sprint(v...)
	logging.Logger.Debug("", zap.String("paho", s))
}

func (l logger) Printf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	logging.Logger.Debug("", zap.String("paho", s))
}
