package config

import (
	"github.com/harluo/config"
)

type Baidu struct {
	Client *Client `default:"{}" json:"client,omitempty" validate:"required"`
}

func newBaidu(config config.Getter) (baidu *Baidu, err error) {
	baidu = new(Baidu)
	err = config.Get(&struct {
		Baidu *Baidu `default:"{}" json:"baidu,omitempty" validate:"required"`
	}{
		Baidu: baidu,
	})

	return
}
