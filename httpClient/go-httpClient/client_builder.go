package apw_http

import (
	"net/http"
	"time"
)

type clientBuilder struct {
	maxIdleConnection int
	connectionTimeout time.Duration
	responseTimeout   time.Duration
	disableTimeout    bool
	headers           http.Header
}

type ClientBuilder interface {
	SetHeaders(headers http.Header) ClientBuilder
	SetConnectionTimeout(timeout time.Duration) ClientBuilder
	SetResponsetimeout(rTimeout time.Duration) ClientBuilder
	SetMaxIdleConnection(i int) ClientBuilder
	DisableTimeout(distable bool) ClientBuilder
	Build() Client
}

func NewClientBuilder() ClientBuilder {
	return &clientBuilder{}
}

func (c *clientBuilder) Build() Client {
	client := &httpClient{
		builder: c,
	}
	return client
}

func (c *clientBuilder) SetHeaders(headers http.Header) ClientBuilder {
	c.headers = headers
	return c
}

func (c *clientBuilder) SetConnectionTimeout(timeout time.Duration) ClientBuilder {
	c.connectionTimeout = timeout
	return c
}

func (c *clientBuilder) SetResponsetimeout(rTimeout time.Duration) ClientBuilder {
	c.responseTimeout = rTimeout
	return c
}

func (c *clientBuilder) SetMaxIdleConnection(i int) ClientBuilder {
	c.maxIdleConnection = i
	return c
}

func (c *clientBuilder) DisableTimeout(distable bool) ClientBuilder {
	c.disableTimeout = distable
	return c
}
