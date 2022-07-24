package client

import (
	"encoding/hex"

	"github.com/emmansun/gmsm/sm9"
)

type UserConfig struct {
	Uid                    string `yaml:"uid"`
	Hid                    byte   `yaml:"hid"`
	EncryptPrivateKey      string `yaml:"encrypt_private_key"`
	SignPrivateKey         string `yaml:"sign_private_key"`
	EncryptMasterPublicKey string `yaml:"encrypt_master_public_key"`
	SignMasterPublicKey    string `yaml:"sign_master_public_key"`
}

type User struct {
	Uid                    []byte
	Hid                    byte
	SessionKey             []byte
	encryptPrivateKey      *sm9.EncryptPrivateKey
	signPrivateKey         *sm9.SignPrivateKey
	encryptMasterPublicKey *sm9.EncryptMasterPublicKey
	signMasterPublicKey    *sm9.SignMasterPublicKey
}

func NewUser(conf *UserConfig) *User {
	if conf == nil {
		return nil
	}

	if conf.EncryptMasterPublicKey == "" || conf.SignMasterPublicKey == "" {
		return nil
	}

	u := &User{Uid: []byte(conf.Uid), Hid: conf.Hid}
	err := u.SetEncryptMasterPublicKey(conf.EncryptMasterPublicKey)
	if err != nil {
		return nil
	}

	err = u.SetSignMasterPublicKey(conf.SignMasterPublicKey)
	if err != nil {
		return nil
	}

	_ = u.SetEncryptPrivateKey(conf.EncryptPrivateKey)
	_ = u.SetSignPrivateKey(conf.SignPrivateKey)

	return u
}

func (u *User) GetEncryptPrivateKey() *sm9.EncryptPrivateKey {
	return u.encryptPrivateKey
}

func (u *User) SetEncryptPrivateKey(encPrivateKey string) error {
	buf, err := hex.DecodeString(encPrivateKey)
	if err != nil {
		return err
	}

	key := new(sm9.EncryptPrivateKey)
	err = key.UnmarshalASN1(buf)
	if err != nil {
		return err
	}

	u.encryptPrivateKey = key
	return nil
}

func (u *User) GetSignPrivateKey() *sm9.SignPrivateKey {
	return u.signPrivateKey
}

func (u *User) SetSignPrivateKey(signPrivateKey string) error {
	buf, err := hex.DecodeString(signPrivateKey)
	if err != nil {
		return err
	}

	key := new(sm9.SignPrivateKey)
	err = key.UnmarshalASN1(buf)
	if err != nil {
		return err
	}

	u.signPrivateKey = key
	return nil
}

func (u *User) GetEncryptMasterPublicKey() *sm9.EncryptMasterPublicKey {
	return u.encryptMasterPublicKey
}

func (u *User) SetEncryptMasterPublicKey(encMasterPublicKey string) error {
	buf, err := hex.DecodeString(encMasterPublicKey)
	if err != nil {
		return err
	}

	key := new(sm9.EncryptMasterPublicKey)
	err = key.UnmarshalASN1(buf)
	if err != nil {
		return err
	}

	u.encryptMasterPublicKey = key
	return nil
}

func (u *User) GetSignMasterPublicKey() *sm9.SignMasterPublicKey {
	return u.signMasterPublicKey
}

func (u *User) SetSignMasterPublicKey(signMasterPublicKey string) error {
	buf, err := hex.DecodeString(signMasterPublicKey)
	if err != nil {
		return err
	}

	key := new(sm9.SignMasterPublicKey)
	err = key.UnmarshalASN1(buf)
	if err != nil {
		return err
	}

	u.signMasterPublicKey = key
	return nil
}
