package core

import (
	"context"
	"fmt"
	"time"

	"github.com/goexl/exception"
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

func (c *Client) Do(ctx context.Context, url string, req any, rsp any) (code uint32, err error) {
	response := new(dto.Response)
	response.Result = rsp

	request := c.http.NewRequest()
	request.SetContext(ctx).SetBody(req).SetResult(response)

	url = fmt.Sprintf("https://aip.baidubce.com/rest/2.0/%s", url)
	fields := gox.Fields[any]{
		field.New("url", url),
	}
	if token, pte := c.pickToken(ctx); nil != pte {
		err = pte
	} else if hpr, hpe := request.SetQueryParam("access_token", token).Post(url); nil != hpe {
		err = hpe
	} else if hpr.IsError() {
		bodyField := field.New("body", string(hpr.Body()))
		message := "百度服务器返回错误"
		err = exception.New().Code(1).Message(message).Field(bodyField, fields...).Build()
		c.logger.Error(message, bodyField, fields...)
	} else if response.IsError() {
		bodyField := field.New("body", string(hpr.Body()))
		message := "接口调用出错"
		c.logger.Warn(message, bodyField, fields...)
	}

	// 将代码回传给上级调用
	code = response.Code

	return
}

func (c *Client) pickToken(ctx context.Context) (token string, err error) {
	now := time.Now()
	if c.expired.Add(-time.Minute).After(now) { // 留一分钟作为缓冲
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
	rsp := new(dto.TokenRsp)
	if hpr, hpe := request.SetContext(ctx).SetQueryParams(params).SetResult(rsp).Post(url); nil != hpe {
		err = hpe
	} else if hpr.IsError() {
		bodyField := field.New("body", string(hpr.Body()))
		err = exception.New().Code(1).Message("百度服务器返回错误").Field(bodyField, fields...).Build()
		c.logger.Error("百度服务器返回错误", bodyField, fields...)
	} else {
		c.token = rsp.Access
		c.expired = time.Now().Add(time.Second * rsp.Duration)
	}

	return
}
