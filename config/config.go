package config

import (
	"fmt"
	"log"

	"github.com/2516319251/boosters/encoding"
)

var _ Config = (*config)(nil)

type Config interface {
	Load() error
	Scan(v interface{}) error
}

type config struct {
	opts   options
	kv     *KeyValue
}

func New(opts ...Option) Config {
	o := options{}

	for _, opt := range opts {
		opt(&o)
	}

	return &config{
		opts:   o,
		kv: &KeyValue{},
	}
}

func (c *config) Load() (err error) {
	c.kv, err = c.opts.source.Load()
	if err != nil {
		return err
	}
	log.Printf("[CONFIG] loaded: %s, format: %s\n", c.kv.Key, c.kv.Format)
	return nil
}

func (c *config) Scan(v interface{}) error {
	if codec := encoding.GetCodec(c.kv.Format); codec != nil {
		return codec.Unmarshal(c.kv.Value, v)
	}
	return fmt.Errorf("unsupported key: %s format: %s", c.kv.Key, c.kv.Format)
}
