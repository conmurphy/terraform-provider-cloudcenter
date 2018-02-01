package main

import (
	"github.com/cloudcenter"
)

type Config struct {
	Username string
	Password string
	Base_url string
}

func (c *Config) Client() *cloudcenter.Client {
	return cloudcenter.NewClient(c.Username, c.Password, c.Base_url)
}
