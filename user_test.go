package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	t.Run("empty config", func(t *testing.T) {
		conf := &UserConfig{}
		u := NewUser(conf)
		assert.Nil(t, u)
	})

	t.Run("empty encrypt and sign private key", func(t *testing.T) {
		conf := &UserConfig{
			Uid:                    "uid",
			Hid:                    1,
			EncryptMasterPublicKey: "034200049174542668e8f14ab273c0945c3690c66e5dd09678b86f734c4350567ed0628354e598c6bf749a3dacc9fffedd9db6866c50457cfc7aa2a4ad65c3168ff74210",
			SignMasterPublicKey:    "03818200049f64080b3084f733e48aff4b41b565011ce0711c5e392cfb0ab1b6791b94c40829dba116152d1f786ce843ed24a3b573414d2177386a92dd8f14d65696ea5e3269850938abea0112b57329f447e3a0cbad3e2fdb1a77f335e89e1408d0ef1c2541e00a53dda532da1a7ce027b7a46f741006e85f5cdff0730e75c05fb4e3216d",
		}
		u := NewUser(conf)

		assert.NotNil(t, u)
		assert.Nil(t, u.GetEncryptPrivateKey())
		assert.Nil(t, u.GetSignPrivateKey())
		assert.NotNil(t, u.GetEncryptMasterPublicKey())
		assert.NotNil(t, u.GetSignMasterPublicKey())
	})

	t.Run("full config", func(t *testing.T) {
		conf := &UserConfig{
			Uid:                    "uid",
			Hid:                    1,
			EncryptPrivateKey:      "038182000461862972d32c7c0fe4df5d7143e9e11c8f429844818501aa877c006ed652496f12b53ae7c707efcbb945a68b41d0b2ac17b6a5d56244ec21e175e9307fe0ba83471fe232f5aa55d24d681789f4507540ccc4d96eb3031de7efd229391759f58636f2e5db70e52f892edb0fbcf98467ca61366e4935564a1013cae7a3db3dc9a8",
			SignPrivateKey:         "0342000449aa82a102d0cb1be45a44eec5e66fbaf289e438f7bf16ce136dbeb252ed17293a7f17e4297501d2310c86324a0a9822537631b6a1a623d55959d726e253ff76",
			EncryptMasterPublicKey: "034200049174542668e8f14ab273c0945c3690c66e5dd09678b86f734c4350567ed0628354e598c6bf749a3dacc9fffedd9db6866c50457cfc7aa2a4ad65c3168ff74210",
			SignMasterPublicKey:    "03818200049f64080b3084f733e48aff4b41b565011ce0711c5e392cfb0ab1b6791b94c40829dba116152d1f786ce843ed24a3b573414d2177386a92dd8f14d65696ea5e3269850938abea0112b57329f447e3a0cbad3e2fdb1a77f335e89e1408d0ef1c2541e00a53dda532da1a7ce027b7a46f741006e85f5cdff0730e75c05fb4e3216d",
		}
		u := NewUser(conf)

		assert.NotNil(t, u)
		assert.NotNil(t, u.GetEncryptPrivateKey())
		assert.NotNil(t, u.GetSignPrivateKey())
		assert.NotNil(t, u.GetEncryptMasterPublicKey())
		assert.NotNil(t, u.GetSignMasterPublicKey())
	})
}
