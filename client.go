package client

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/url"
	"time"

	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
)

type ClientConfig struct {
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

type Client struct {
	ServerUrl   *url.URL
	AuthHandler *Sm9Auth
	User        *User
	Cm          *autopaho.ConnectionManager
	handler     *handler
	Cancel      context.CancelFunc
	Config      *ClientConfig
}

func (c *Client) Connect() error {
	// Create a handler that will deal with incoming messages
	c.handler = NewHandler(c.Config.WriteToDisk, c.Config.OutputFileName, c.Config.WriteToStdOut)

	cliCfg := autopaho.ClientConfig{
		BrokerUrls:        []*url.URL{c.ServerUrl},
		KeepAlive:         c.Config.Keepalive,
		ConnectRetryDelay: time.Duration(c.Config.ConnectRetryDelay) * time.Millisecond,
		OnConnectionUp: func(cm *autopaho.ConnectionManager, connAck *paho.Connack) {
			fmt.Println("mqtt connection up")
			if _, err := cm.Subscribe(context.Background(), &paho.Subscribe{
				Subscriptions: map[string]paho.SubscribeOptions{
					c.Config.Topic: {QoS: c.Config.Qos},
				},
			}); err != nil {
				fmt.Printf("failed to subscribe (%s). This is likely to mean no messages will be received.", err)
				return
			}
			fmt.Println("mqtt subscription made")
		},
		OnConnectError: func(err error) { fmt.Printf("error whilst attempting connection: %s\n", err) },
		ClientConfig: paho.ClientConfig{
			ClientID: c.Config.ClientID,
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
					Value: c.Config.ClientName,
				},
			},
		}

		return connect
	})

	if c.Config.Debug {
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
	subPacket := &paho.Subscribe{
		Subscriptions: map[string]paho.SubscribeOptions{
			topic: {QoS: c.Config.Qos},
		},
	}
	if _, err := c.Cm.Subscribe(context.Background(), subPacket); err != nil {
		fmt.Printf("failed to subscribe (%s). This is likely to mean no messages will be received.", err)
		return err
	}

	fmt.Printf("Subscribed to %s", topic)
	return nil
}

func (c *Client) Publish(topic, payload string) error {
	pubPacket := &paho.Publish{
		Topic:   topic,
		QoS:     byte(0),
		Payload: []byte(payload),
	}

	if _, err := c.Cm.Publish(context.Background(), pubPacket); err != nil {
		fmt.Println(err)
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
