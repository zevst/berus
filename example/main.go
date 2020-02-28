package main

import (
	"bytes"
	"context"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"github.com/zevst/berus"
	"github.com/zevst/zlog"
	"go.uber.org/zap"
	"net/http"
)

var yamlExample = []byte(`
system:
 transport:
   _type: rest
   http_client:
     timeout: 30s
`)

var Config struct {
	System *System `mapstructure:"system"`
}

type System struct {
	Transport berus.Custom `mapstructure:"transport"`
}

type Rest struct {
	Client *http.Client `mapstructure:"http_client"`
}

func (d *Rest) Send(context.Context, *http.Request) (*http.Response, error) {
	// Some kind of logic ...
	return &http.Response{}, nil
}

func init() {
	berus.RegisterCustom("rest", &Rest{})

	// This is to simplify the code.
	viper.SetConfigType("yaml")
	if err := viper.ReadConfig(bytes.NewBuffer(yamlExample)); err != nil {
		zlog.Fatal(err.Error())
	}
	err := viper.Unmarshal(&Config, viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToSliceHookFunc(","),
		// Custom Decode Hook Function
		berus.CustomHookFunc,
	)))
	if err != nil {
		zlog.Fatal(err.Error())
	}
}

func main() {
	switch protocol := Config.System.Transport.(type) {
	case *Rest:
		response, err := protocol.Send(context.Background(), &http.Request{})
		if err != nil {
			zlog.Fatal(err.Error())
			return // zlog.Fatal method calls os.Exit, but code analyzer doesn't know it :(
		}
		zlog.Info("response",
			zap.Int("http_code", response.StatusCode),
			zap.String("http_status", response.Status),
		)
	default:
		zlog.Fatal("Unsupported protocol", zap.Any("protocol", protocol))
	}
}
