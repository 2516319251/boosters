package main

import (
	"flag"
	"net/url"

	"github.com/2516319251/boosters"
	"github.com/2516319251/boosters/config"
	"github.com/2516319251/boosters/config/file"
	_ "github.com/2516319251/boosters/encoding/yaml"
	"github.com/2516319251/boosters/examples/greeter/conf"
	"github.com/2516319251/boosters/examples/greeter/data"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// flagconf is the config flag.
	flagconf string
)

func init() {
	flag.StringVar(&flagconf, "conf", "./config.yaml", "config path, eg: -conf config.yaml")
}

func main() {
	flag.Parse()
	c := config.New(config.WithSource(file.Load(flagconf)))
	if e := c.Load(); e != nil {
		panic(e)
	}

	var cfg conf.Config
	if e := c.Scan(&cfg); e != nil {
		panic(e)
	}

	client := data.NewClient()
	defer client.Close()

	u, e := url.Parse("grpc://127.0.0.1:9100?isSecure=true")
	if e != nil {
		panic(e)
	}

	u1, e := url.Parse("grpc://127.0.0.1:9200?isSecure=true")
	if e != nil {
		panic(e)
	}

	opts := []boosters.Option{
		boosters.ID(""),
		boosters.Name("greeter"),
		boosters.Version("v1.0.0"),
		boosters.Endpoint(u, u1),
		boosters.Registrar(data.NewRegistry(client)),
		boosters.Server(NewServer(cfg.Server.Grpc)),
	}

	app := boosters.New(opts...)
	if err := app.Run(); err != nil {
		panic(err)
	}
}
