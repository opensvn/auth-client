package client

import (
	"encoding/hex"
	"github.com/emmansun/gmsm/sm9"
)

type User struct {
	Uid               []byte
	Hid               byte
	encryptPrivateKey *sm9.EncryptPrivateKey
	signPrivateKey    *sm9.SignPrivateKey
}

func (u *User) GetEncryptPrivateKey() *sm9.EncryptPrivateKey {
	return u.encryptPrivateKey
}

func (u *User) SetEncryptPrivateKey(s string) error {
	buf, err := hex.DecodeString(s)
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

func (u *User) SetSignPrivateKey(s string) error {
	buf, err := hex.DecodeString(s)
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
	return &u.encryptPrivateKey.EncryptMasterPublicKey
}

func (u *User) SetEncryptMasterPublicKey(s string) error {
	buf, err := hex.DecodeString(s)
	if err != nil {
		return err
	}

	key := new(sm9.EncryptMasterPublicKey)
	err = key.UnmarshalASN1(buf)
	if err != nil {
		return err
	}

	u.encryptPrivateKey.EncryptMasterPublicKey = *key
	return nil
}

func (u *User) GetSignMasterPublicKey() *sm9.SignMasterPublicKey {
	return &u.signPrivateKey.SignMasterPublicKey
}

func (u *User) SetSignMasterPublicKey(s string) error {
	buf, err := hex.DecodeString(s)
	if err != nil {
		return err
	}

	key := new(sm9.SignMasterPublicKey)
	err = key.UnmarshalASN1(buf)
	if err != nil {
		return err
	}

	u.signPrivateKey.SignMasterPublicKey = *key
	return nil
}
