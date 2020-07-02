package locker_client

import (
	"os"
	"sync"
	"testing"

	"github.com/coinexchain/cloudlocker"
)

var testUrl4Server = "127.0.0.1:33300"
var testUrl = "http://" + testUrl4Server
var testPath = "./tmp"

func TestMultiGoroutineWrite(_ *testing.T) {
	wg := sync.WaitGroup{}
	db := CloudLockerClient{
		url: testUrl,
	}
	server, _ := cloudlocker.NewLockerServer(testPath, testUrl4Server)
	go server.Start()
	defer server.Stop()

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(index int) {
			for j := 0; j < 1000; j++ {
				key := []byte{byte(index)}
				value := []byte{byte(j % 100)}
				_ = db.Set(key, value)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	for i := 0; i < 100; i++ {
		v, _ := db.Get([]byte{byte(i)})
		if len(v) == 0 || v[0] != 99 {
			panic("v should be 99")
		}
	}

	_ = os.RemoveAll(testPath)
}

func TestBasics(t *testing.T) {
	db := CloudLockerClient{
		url: testUrl,
	}

	server, _ := cloudlocker.NewLockerServer(testPath, testUrl4Server)
	go server.Start()
	defer server.Stop()

	for i := 0; i < 100; i++ {
		_ = db.Set([]byte{byte(i)}, []byte{byte(i + 1)})
		v, _ := db.Get([]byte{byte(i)})
		if len(v) == 0 || v[0] != byte(i+1) {
			panic("value not set exactly")
		}
	}

	_ = os.RemoveAll(testPath)
}
