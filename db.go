package locker_client

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type DB interface {
	Get(key []byte) ([]byte, error)
	Set(key, value []byte) error
	Delete(key []byte) error
}

var _ DB = CloudLockerClient{}

type CloudLockerClient struct {
	url string
}

func NewCloudLockerClient(url string) *CloudLockerClient {
	return &CloudLockerClient{url: url}
}

//if key is not exist, return nil slice and error is nil
func (c CloudLockerClient) Get(key []byte) ([]byte, error) {
	resp, _ := http.Post(c.url+"/get", "application/json", strings.NewReader(string(key)))
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}

func (c CloudLockerClient) Set(key, value []byte) error {
	e := entry{
		K: string(key),
		V: string(value),
	}
	b, _ := json.Marshal(e)
	_, err := http.Post(c.url+"/set", "application/json", strings.NewReader(string(b)))
	return err
}

func (c CloudLockerClient) Delete(key []byte) error {
	_, err := http.Post(c.url+"/delete", "application/json", strings.NewReader(string(key)))
	return err
}

type entry struct {
	K string `json:"k"`
	V string `json:"v"`
}
