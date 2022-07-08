package client

import (
	"encoding/hex"

	"github.com/emmansun/gmsm/sm9"
)

type User struct {
	Uid               []byte
	Hid               byte
	EncryptPrivateKey *sm9.EncryptPrivateKey
	SignPrivateKey    *sm9.SignPrivateKey
}

func (u *User) GetEncryptPrivateKey() *sm9.EncryptPrivateKey {
	return u.EncryptPrivateKey
}

func (u *User) SetEncryptPrivateKey(s string) error {
	buf, err := hex.DecodeString(s)
	if err != nil {
		return err
	}

	err = u.EncryptPrivateKey.UnmarshalASN1(buf)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) GetSignPrivateKey() *sm9.SignPrivateKey {
	return u.SignPrivateKey
}

func (u *User) SetSignPrivateKey(s string) error {
	buf, err := hex.DecodeString(s)
	if err != nil {
		return err
	}

	err = u.SignPrivateKey.UnmarshalASN1(buf)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) GetEncryptMasterPublicKey() *sm9.EncryptMasterPublicKey {
	return &u.EncryptPrivateKey.EncryptMasterPublicKey
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

	u.EncryptPrivateKey.EncryptMasterPublicKey = *key
	return nil
}

func (u *User) GetSignMasterPublicKey() *sm9.SignMasterPublicKey {
	return &u.SignPrivateKey.SignMasterPublicKey
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

	u.SignPrivateKey.SignMasterPublicKey = *key
	return nil
}
