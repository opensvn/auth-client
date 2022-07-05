package client

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
)

type Client struct {
	c           *paho.Client
	ServerUrl   *url.URL
	AuthHandler *Sm9Auth
	User        *User
	Cm          *autopaho.ConnectionManager
	handler     *handler
	Cancel      context.CancelFunc

	ClientID          string
	ClientName        string
	Topic             string
	Qos               byte
	Keepalive         uint16
	ConnectRetryDelay uint16
	WriteToStdOut     bool
	WriteToDisk       bool
	OutputFileName    string
	Debug             bool
}

func (c *Client) Connect() error {
	// Create a handler that will deal with incoming messages
	c.handler = NewHandler(c.WriteToDisk, c.OutputFileName, c.WriteToStdOut)

	cliCfg := autopaho.ClientConfig{
		BrokerUrls:        []*url.URL{c.ServerUrl},
		KeepAlive:         c.Keepalive,
		ConnectRetryDelay: time.Duration(c.ConnectRetryDelay) * time.Millisecond,
		OnConnectionUp: func(cm *autopaho.ConnectionManager, connAck *paho.Connack) {
			fmt.Println("mqtt connection up")
			if _, err := cm.Subscribe(context.Background(), &paho.Subscribe{
				Subscriptions: map[string]paho.SubscribeOptions{
					c.Topic: {QoS: c.Qos},
				},
			}); err != nil {
				fmt.Printf("failed to subscribe (%s). This is likely to mean no messages will be received.", err)
				return
			}
			fmt.Println("mqtt subscription made")
		},
		OnConnectError: func(err error) { fmt.Printf("error whilst attempting connection: %s\n", err) },
		ClientConfig: paho.ClientConfig{
			ClientID: c.ClientID,
			Router: paho.NewSingleHandlerRouter(func(m *paho.Publish) {
				c.handler.handle(m)
			}),
			OnClientError: func(err error) { fmt.Printf("server requested disconnect: %s\n", err) },
			OnServerDisconnect: func(d *paho.Disconnect) {
				if d.Properties != nil {
					fmt.Printf("server requested disconnect: %s\n", d.Properties.ReasonString)
				} else {
					fmt.Printf("server requested disconnect; reason code: %d\n", d.ReasonCode)
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
					Value: c.ClientName,
				},
			},
		}

		return connect
	})

	if c.Debug {
		cliCfg.Debug = logger{prefix: "autoPaho"}
		cliCfg.PahoDebug = logger{prefix: "paho"}
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
	if _, err := c.c.Subscribe(context.Background(), &paho.Subscribe{
		Subscriptions: map[string]paho.SubscribeOptions{
			topic: {QoS: byte(0), NoLocal: true},
		},
	}); err != nil {
		log.Println(err)
		return err
	}

	log.Printf("Subscribed to %s", topic)
	return nil
}

func (c *Client) Publish(topic, payload string) error {
	pb := &paho.Publish{
		Topic:   topic,
		QoS:     byte(0),
		Payload: []byte(payload),
	}

	if _, err := c.c.Publish(context.Background(), pb); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (c *Client) Disconnect() error {
	defer c.handler.Close()
	defer c.Cancel()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return c.Cm.Disconnect(ctx)
}

// logger implements the paho.Logger interface
type logger struct {
	prefix string
}

// Println is the library provided NOOPLogger's
// implementation of the required interface function()
func (l logger) Println(v ...interface{}) {
	fmt.Println(append([]interface{}{l.prefix + ":"}, v...)...)
}

// Printf is the library provided NOOPLogger's
// implementation of the required interface function(){}
func (l logger) Printf(format string, v ...interface{}) {
	if len(format) > 0 && format[len(format)-1] != '\n' {
		format = format + "\n" // some log calls in paho do not add \n
	}
	fmt.Printf(l.prefix+":"+format, v...)
}
