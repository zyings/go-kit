package config

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func CreateFile() {
	fp, _ := os.Create("test.json")
	_, _ = fp.WriteString(`{
  "Host": "localhost",
  "Port": 6060,
  "Log": [{
    "File": "test.info",
    "MaxAge": "24h",
    "Format": "json"
  }, {
    "File": "test.warn",
    "MaxAge": "24h",
    "Format": "text"
  }]
}`)
	_ = fp.Close()
}

func DeleteFile() {
	_ = os.Remove("test.json")
}

func TestConfigExample1(t *testing.T) {
	Convey("TestConfigExample1", t, func() {
		CreateFile()
		//provider, _ := NewLocalProvider("test.json")
		//conf, err := NewConfig(&Json5Decoder{}, provider, nil)

		conf, err := NewConfigWithBaseFile("testfile/base.json")
		So(err, ShouldBeNil)
		//So(conf.GetInt("Port"), ShouldEqual, 6060)
		//So(conf.GetString("Host"), ShouldEqual, "localhost")
		//So(conf.GetString("Logger[0].File"), ShouldEqual, "test.info")
		//So(conf.GetDuration("Logger[0].MaxAge"), ShouldEqual, 24*time.Hour)
		//So(conf.GetString("Logger[1].File"), ShouldEqual, "test.warn")
		//So(conf.GetDuration("Logger[1].MaxAge"), ShouldEqual, 168*time.Hour)

		if err := conf.Watch(); err != nil {
			fmt.Println(err)
		}

		opt1 := &LogOption{}
		So(conf.Sub("Logger[0]").Unmarshal(opt1), ShouldBeNil)
		fmt.Println(opt1)
		fmt.Println(reflect.TypeOf(opt1))
		fmt.Println(conf.Get("OSS"))

		opt := conf.Bind("Logger[0]", LogOption{}, OnSucc(func(c *Config) {
			fmt.Println("update logger succ", c.GetString("Logger[0].Name"))
		}), OnFail(func(err error) {
			fmt.Println(err)

		}))
		port := conf.Int("Port", OnSucc(func(c *Config) {
			fmt.Println("update port success")
		}))
		fmt.Println(conf.Get("Logger[0]"))
		for i := 0; i < 60; i++ {
			fmt.Println(port.Get())
			fmt.Println(opt.Load().(LogOption))
			//time.Sleep(time.Second)
		}

		DeleteFile()
	})
}

type LogOption struct {
	File   string
	Name   string
	MaxAge time.Duration
}

type MyOption struct {
	Host   string
	Port   int
	Logger []*LogOption
}

func TestConfigExample2(t *testing.T) {

	Convey("TestConfigExample2", t, func() {
		CreateFile()
		//provider, _ := NewLocalProvider("testfile/test.json")
		//abc, _ := base64.StdEncoding.DecodeString("IrjXy4vx7iwgCLaUeu5TVUA9TkgMwSw3QWcgE/IW5W0=")
		//cipher, _ := NewAESCipher(abc)
		//conf, err := NewConfig(&Json5Decoder{}, provider, cipher)
		conf, err := NewConfigWithBaseFile("testfile/base.json")
		So(err, ShouldBeNil)

		buf, _ := json.MarshalIndent(conf.storage.root, "  ", "  ")
		fmt.Println(string(buf))

		opt := &MyOption{}
		So(conf.Unmarshal(opt), ShouldBeNil)
		So(opt.Host, ShouldEqual, "localhost")
		fmt.Println(opt)
		fmt.Println(opt.Logger[0])
		DeleteFile()
	})
}
