package client

import (
	"net/http"
	"os"
	"pingrobot/logger"
	"time"
)

type Client struct {
	Client http.Client
}

func NewClient() *Client {
	return &Client{
		Client: http.Client{
			Timeout: time.Second * 3,
			Transport: &logger.LoggingRoundTripper{
				Logger: os.Stdout,
				Next:   http.DefaultTransport,
			},
		},
	}
}
