package main

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Ssh         string `json:"ssh"`
	Source      string `json:"source"`
	IsEncrypted bool   `json:"isEncrypted"`
	Password    string `json:"password"`
}

func NewConfig(filename string) (config *Config, err error) {
	config = &Config{}
	var data []byte
	data, err = ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, config)
	if err != nil {
		return
	}
	return
}

func (c *Config) PasswordFile() string {
	if c.Password == "" {
		return "passwd"
	}
	return c.Password
}
