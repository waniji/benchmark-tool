package main

import (
	"encoding/json"
	"errors"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"strings"
)

type Config struct {
	URLs          []string `json:"url"`
	MaxAccess     int      `json:"count"`
	MaxWorkers    int      `json:"worker"`
	BasicAuthUser string   `json:"basic-auth-user"`
	BasicAuthPass string   `json:"basic-auth-pass"`
	Format        string   `json:"format"`
	Sort          string   `json:"sort"`
}

func (config *Config) Create(context *cli.Context) error {

	configFile := context.String("config-file")
	if configFile != "" {
		if err := config.setFromConfigFile(configFile); err != nil {
			return err
		}
	}

	config.setFromArgs(context)

	if err := config.validate(); err != nil {
		return err
	}
	return nil
}

func (c *Config) setFromConfigFile(configFile string) error {
	configData, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(configData, c)
	if err != nil {
		return err
	}

	return nil
}

func (config *Config) setFromArgs(context *cli.Context) {
	if context.String("url") != "" {
		config.URLs = strings.Split(context.String("url"), ",")
	}
	if context.Int("count") != 0 {
		config.MaxAccess = context.Int("count")
	}
	if context.Int("worker") != 0 {
		config.MaxWorkers = context.Int("worker")
	}
	if context.String("basic-auth-user") != "" {
		config.BasicAuthUser = context.String("basic-auth-user")
	}
	if context.String("basic-auth-pass") != "" {
		config.BasicAuthPass = context.String("basic-auth-pass")
	}
	if context.String("format") != "" {
		config.Format = context.String("format")
	}
	if context.String("sort") != "" {
		config.Sort = context.String("sort")
	}
}

func (c *Config) validate() error {
	if len(c.URLs) == 0 {
		return errors.New("urlは必須です")
	}
	if c.MaxAccess == 0 {
		return errors.New("countは必須です")
	}
	if c.MaxWorkers == 0 {
		return errors.New("workerは必須です")
	}
	return nil
}
