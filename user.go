package client

import (
	"encoding/hex"
	"github.com/emmansun/gmsm/sm9"
	"github.com/opensvn/kgc"
)

var Kgc *kgc.Kgc
var CurrentUser *User

func init() {
	k, err := kgc.New()
	if err != nil {
		panic(err)
	}
	Kgc = k
}

type User struct {
	Uid               []byte
	Hid               byte
	signPrivateKey    *sm9.SignPrivateKey
	signPublicKey     *sm9.SignMasterPublicKey
	encryptPrivateKey *sm9.EncryptPrivateKey
	encryptPublicKey  *sm9.EncryptMasterPublicKey
}

func NewUser(kgc *kgc.Kgc, uid []byte, hid byte) *User {
	signKey, err := kgc.GenerateUserSignKey(uid, hid)
	if err != nil {
		return nil
	}

	encryptKey, err := kgc.GenerateUserEncryptKey(uid, hid)
	if err != nil {
		return nil
	}

	return &User{
		Uid:               uid,
		Hid:               hid,
		signPrivateKey:    signKey,
		signPublicKey:     signKey.MasterPublic(),
		encryptPrivateKey: encryptKey,
		encryptPublicKey:  encryptKey.MasterPublic(),
	}
}

func (u *User) GetSignMasterPublicKeyASN1() string {
	buf, err := u.signPublicKey.MarshalASN1()
	if err != nil {
		return ""
	}

	return hex.EncodeToString(buf)
}

func (u *User) SetSignMasterPublicKeyASN1(s string) error {
	buf, err := hex.DecodeString(s)
	if err != nil {
		return err
	}

	masterKey := new(sm9.SignMasterPublicKey)
	err = masterKey.UnmarshalASN1(buf)
	if err != nil {
		return err
	}

	u.signPublicKey = masterKey
	return nil
}

func (u *User) GetEncryptMasterPublicKeyASN1() string {
	buf, err := u.encryptPublicKey.MarshalASN1()
	if err != nil {
		return ""
	}

	return hex.EncodeToString(buf)
}

func (u *User) SetEncryptMasterPublicKeyASN1(s string) error {
	buf, err := hex.DecodeString(s)
	if err != nil {
		return err
	}

	masterKey := new(sm9.EncryptMasterPublicKey)
	err = masterKey.UnmarshalASN1(buf)
	if err != nil {
		return err
	}

	u.encryptPublicKey = masterKey
	return nil
}
