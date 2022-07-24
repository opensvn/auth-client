package main

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/emmansun/gmsm/sm9"
	"github.com/opensvn/auth-client"
	"github.com/opensvn/auth-client/cmd/config"
)

type RegisterRequest struct {
	Id         string `json:"id"`
	Username   string `json:"username"`
	Eid        string `json:"eid"`
	Random     []byte `json:"random"`
	DeviceType string `json:"device_type"`
}

type Keys struct {
	SignKey    string `json:"signkey"`
	EncryptKey string `json:"encryptkey"`
}

type KeyResponse struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
	Data Keys   `json:"data"`
}

func ApplyKey(conf *config.Config, user *client.User) ([]byte, error) {
	buf, err := getRandom(16)
	if err != nil {
		panic(err)
	}

	random := buf[:]

	random1, err := sm9.EncryptASN1(rand.Reader, user.GetEncryptMasterPublicKey(), []byte("pkg"), 1, buf)
	if err != nil {
		return nil, err
	}

	req := RegisterRequest{
		Id:       conf.User.Uid,
		Username: conf.Mqtt.ClientName,
		Eid:      conf.User.Uid,
		Random:   random1,
		DeviceType: conf.Mqtt.DeviceType,
	}

	buf, err = json.Marshal(req)
	if err != nil {
		return nil, err
	}

	url := "http://" + conf.Addr.Ra + "/register"
	_, err = Post(url, buf)
	if err != nil {
		return nil, err
	}

	return random, nil
}

func queryKey(conf *config.Config) (*Keys, error) {
	url := "http://" + conf.Addr.Platform + "/identificationinfo/identificationinfo/keys?id=" + conf.User.Uid
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var resp KeyResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

func getRandom(len int) ([]byte, error) {
	buf := make([]byte, len)
	_, err := rand.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func Post(url string, body []byte) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second} // 设置请求超时时长5s
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}
