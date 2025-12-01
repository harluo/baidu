package core

import (
	"context"
	"fmt"
	"time"

	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
	"github.com/goexl/http"
	"github.com/goexl/log"
	"github.com/harluo/baidu/internal/config"
	"github.com/harluo/baidu/internal/core/internal/dto"
)

type Client struct {
	token   string
	expired time.Time

	config *config.Client
	http   *http.Client
	logger log.Logger
}

func newClient(
	config *config.Client,

	http *http.Client, logger log.Logger,
) *Client {
	return &Client{
		config: config,

		http:   http,
		logger: logger,
	}
}

func (c *Client) Send(ctx context.Context, url string, req any, rsp any) (err error) {
	url = fmt.Sprintf("https://aip.baidubce.com/rest/2.0/%s", url)
	request := c.http.NewRequest()
	request.SetContext(ctx).SetBody(req).SetResult(rsp)

	fields := gox.Fields[any]{
		field.New("url", url),
	}
	if token, pte := c.pickToken(ctx); nil != pte {
		err = pte
	} else if hpr, hpe := request.SetQueryParam("access_token", token).Post(url); nil != hpe {
		err = hpe
	} else if hpr.IsError() {
		c.logger.Error("百度服务器返回错误", field.New("status", hpr.Status()), fields...)
	}

	return
}

func (c *Client) pickToken(ctx context.Context) (token string, err error) {
	now := time.Now()
	if now.After(c.expired.Add(time.Minute)) { // 留一分钟作为缓冲
		token = c.token
	} else {
		token, err = c.getToken(ctx)
	}

	return
}

func (c *Client) getToken(ctx context.Context) (token string, err error) {
	url := "https://aip.baidubce.com/oauth/2.0/token"
	request := c.http.NewRequest()
	params := map[string]string{
		"grant_type":    "client_credentials",
		"client_id":     c.config.Id,
		"client_secret": c.config.Secret,
	}

	fields := gox.Fields[any]{
		field.New("url", url),
	}
	rsp := new(dto.AccessRsp)
	if hpr, hpe := request.SetContext(ctx).SetQueryParams(params).SetResult(rsp).Post(url); nil != hpe {
		err = hpe
	} else if hpr.IsError() {
		c.logger.Error("百度服务器返回错误", field.New("status", hpr.Status()), fields...)
	} else {
		c.token = rsp.Token
		c.expired = rsp.Expired
	}

	return
}
