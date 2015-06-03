package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	"strings"
)

func bench(c *cli.Context) {

	url := c.String("url")
	maxAccess := c.Int("count")
	maxWorkers := c.Int("worker")
	basicAuthUser := c.String("basic-auth-user")
	basicAuthPass := c.String("basic-auth-pass")

	if url == "" {
		fmt.Println("urlは必須です")
		cli.ShowAppHelp(c)
		os.Exit(1)
	}

	fmt.Printf("Total Access Count: %d\n", maxAccess)
	fmt.Printf("Concurrency: %d\n", maxWorkers)
	fmt.Println("--------------------------------------------------")

	manager := &WorkerManager{
		urls:          strings.Split(url, ","),
		basicAuthUser: basicAuthUser,
		basicAuthPass: basicAuthPass,
		maxAccess:     maxAccess,
		maxWorkers:    maxWorkers,
		activeWorker:  make(chan struct{}, maxWorkers),
		doneWorker:    make(chan struct{}, maxAccess),
	}
	manager.Start()

	fmt.Println("--------------------------------------------------")
	fmt.Printf("Total Time           : %d msec\n", manager.result.totalElapsedMsec)
	fmt.Printf("Average Response Time: %d msec\n", manager.result.averageElapsedMsec)
	fmt.Printf("Minimum Response Time: %d msec\n", manager.result.minimumElapsedMsec)
	fmt.Printf("Maximum Response Time: %d msec\n", manager.result.maximumElapsedMsec)
}
