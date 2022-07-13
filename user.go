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
	Uid               []byte
	Hid               byte
	EncryptPrivateKey *sm9.EncryptPrivateKey
	SignPrivateKey    *sm9.SignPrivateKey
}

func NewUser(conf *UserConfig) *User {
	if conf == nil {
		return nil
	}

	if conf.EncryptPrivateKey == "" || conf.EncryptMasterPublicKey == "" ||
		conf.SignPrivateKey == "" || conf.SignMasterPublicKey == "" {
		return nil
	}

	u := &User{Uid: []byte(conf.Uid), Hid: conf.Hid}
	err := u.SetEncryptPrivateKey(conf.EncryptPrivateKey, conf.EncryptMasterPublicKey)
	if err != nil {
		return nil
	}

	err = u.SetSignPrivateKey(conf.SignPrivateKey, conf.SignMasterPublicKey)
	if err != nil {
		return nil
	}

	return u
}

func (u *User) GetEncryptPrivateKey() *sm9.EncryptPrivateKey {
	return u.EncryptPrivateKey
}

func (u *User) SetEncryptPrivateKey(encPrivateKey, encMasterPublicKey string) error {
	buf, err := hex.DecodeString(encPrivateKey)
	if err != nil {
		return err
	}

	u.EncryptPrivateKey = new(sm9.EncryptPrivateKey)
	err = u.EncryptPrivateKey.UnmarshalASN1(buf)
	if err != nil {
		return err
	}

	buf, err = hex.DecodeString(encMasterPublicKey)
	err = u.EncryptPrivateKey.EncryptMasterPublicKey.UnmarshalASN1(buf)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) GetSignPrivateKey() *sm9.SignPrivateKey {
	return u.SignPrivateKey
}

func (u *User) SetSignPrivateKey(signPrivateKey, signMasterPublicKey string) error {
	buf, err := hex.DecodeString(signPrivateKey)
	if err != nil {
		return err
	}

	u.SignPrivateKey = new(sm9.SignPrivateKey)
	err = u.SignPrivateKey.UnmarshalASN1(buf)
	if err != nil {
		return err
	}

	buf, err = hex.DecodeString(signMasterPublicKey)
	err = u.SignPrivateKey.SignMasterPublicKey.UnmarshalASN1(buf)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) GetEncryptMasterPublicKey() *sm9.EncryptMasterPublicKey {
	if u.EncryptPrivateKey == nil {
		return nil
	}
	return &u.EncryptPrivateKey.EncryptMasterPublicKey
}

func (u *User) GetSignMasterPublicKey() *sm9.SignMasterPublicKey {
	if u.SignPrivateKey == nil {
		return nil
	}
	return &u.SignPrivateKey.SignMasterPublicKey
}
