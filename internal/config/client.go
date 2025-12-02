package config

type Client struct {
	// 授权标识
	Key string `json:"key,omitempty" validate:"required"`
	// 授权密码
	Secret string `json:"secret,omitempty" validate:"required"`
}

func newClient(baidu *Baidu) *Client {
	return baidu.Client
}
