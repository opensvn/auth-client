package client

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/eclipse/paho.golang/packets"
	"github.com/eclipse/paho.golang/paho"
	"github.com/emmansun/gmsm/sm9"
)

type Sm9Auth struct {
	Random1 string
	Server  *User
	client *Client
}

func NewSm9Auth(c *Client) *Sm9Auth {
	return &Sm9Auth{client: c}
}

func (s *Sm9Auth) Authenticate(a *paho.Auth) *paho.Auth {
	reauth := &paho.Auth{
		Properties: &paho.AuthProperties{
			AuthMethod: "sm9",
		},
		ReasonCode: packets.AuthReauthenticate,
	}

	s.Server = &User{}
	s.Server.Uid = []byte(a.Properties.User.Get("uid"))
	buf, err := hex.DecodeString(a.Properties.User.Get("hid"))
	if err != nil {
		return reauth
	}
	s.Server.Hid = buf[0]

	buf, err = hex.DecodeString(string(a.Properties.AuthData))
	if err != nil {
		return reauth
	}

	decrypted, err := sm9.DecryptASN1(s.client.User.encryptPrivateKey, s.client.User.Uid, buf)
	if err != nil {
		return reauth
	}

	random1 := decrypted[:len(decrypted)/2]
	if s.Random1 != hex.EncodeToString(random1) {
		return reauth
	}

	random2 := decrypted[len(decrypted)/2:]
	buf, err = sm9.EncryptASN1(rand.Reader, s.client.User.GetEncryptMasterPublicKey(), s.Server.Uid, s.Server.Hid, random2)
	if err != nil {
		return reauth
	}

	return &paho.Auth{
		Properties: &paho.AuthProperties{
			AuthMethod: "sm9",
			AuthData:   []byte(hex.EncodeToString(buf)),
		},
		ReasonCode: packets.AuthContinueAuthentication,
	}
}

func (s *Sm9Auth) Authenticated() {}

func (s *Sm9Auth) GetRandom1(expectedLen int) string {
	buf := make([]byte, expectedLen)
	_, err := rand.Read(buf)
	if err != nil {
		return ""
	}

	s.Random1 = hex.EncodeToString(buf)
	return s.Random1
}
