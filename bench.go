package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"os"
	"strings"
)

type Config struct {
	URLs          []string `json:"url"`
	MaxAccess     int      `json:"count"`
	MaxWorkers    int      `json:"worker"`
	BasicAuthUser string   `json:"basic-auth-user"`
	BasicAuthPass string   `json:"basic-auth-pass"`
}

func (c *Config) Validate() error {
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

func bench(c *cli.Context) {

	var config Config

	configFile := c.String("config-file")
	if configFile != "" {
		configData, err := ioutil.ReadFile(configFile)
		if err != nil {
			fmt.Println(err)
			cli.ShowAppHelp(c)
			os.Exit(1)
		}

		err = json.Unmarshal(configData, &config)
		if err != nil {
			fmt.Println(err)
			cli.ShowAppHelp(c)
			os.Exit(1)
		}
	}

	if c.String("url") != "" {
		config.URLs = strings.Split(c.String("url"), ",")
	}
	if c.Int("count") != 0 {
		config.MaxAccess = c.Int("count")
	}
	if c.Int("worker") != 0 {
		config.MaxWorkers = c.Int("worker")
	}
	if c.String("basic-auth-user") != "" {
		config.BasicAuthUser = c.String("basic-auth-user")
	}
	if c.String("basic-auth-pass") != "" {
		config.BasicAuthPass = c.String("basic-auth-pass")
	}

	if err := config.Validate(); err != nil {
		fmt.Println(err)
		cli.ShowAppHelp(c)
		os.Exit(1)
	}

	fmt.Printf("Total Access Count: %d\n", config.MaxAccess)
	fmt.Printf("Concurrency: %d\n", config.MaxWorkers)
	fmt.Println("--------------------------------------------------")

	manager := &WorkerManager{
		urls:          config.URLs,
		basicAuthUser: config.BasicAuthUser,
		basicAuthPass: config.BasicAuthPass,
		maxAccess:     config.MaxAccess,
		maxWorkers:    config.MaxWorkers,
		activeWorker:  make(chan struct{}, config.MaxWorkers),
		doneWorker:    make(chan struct{}, config.MaxAccess),
	}
	manager.Start()

	fmt.Println("--------------------------------------------------")
	fmt.Printf("Total Time           : %d msec\n", manager.result.totalElapsedMsec)
	fmt.Printf("Average Response Time: %d msec\n", manager.result.averageElapsedMsec)
	fmt.Printf("Minimum Response Time: %d msec\n", manager.result.minimumElapsedMsec)
	fmt.Printf("Maximum Response Time: %d msec\n", manager.result.maximumElapsedMsec)
}
