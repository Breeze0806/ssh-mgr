package dao

import (
	"net/url"
)

type SshConnInfo struct {
	Name     string `json:"-"`
	Group    string `json:"-"`
	Address  string `json:"address"`
	User     string `json:"user"`
	Password string `json:"password"`
}

func (c *SshConnInfo) SftpURL() string {
	URL := &url.URL{
		Scheme: "sftp",
		User:   url.UserPassword(c.User, c.Password),
		Host:   c.Address,
	}

	return URL.String()
}

func (c *SshConnInfo) SshURL() string {
	URL := &url.URL{
		Scheme: "",
		User:   url.User(c.User),
		Host:   c.Address,
	}

	return URL.String()[2:]
}
