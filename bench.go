package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
)

func bench(c *cli.Context) {

	var config Config

	if err := config.Create(c); err != nil {
		fmt.Println(err)
		cli.ShowAppHelp(c)
		os.Exit(1)
	}

	manager := &WorkerManager{
		urls:          config.URLs,
		basicAuthUser: config.BasicAuthUser,
		basicAuthPass: config.BasicAuthPass,
		maxAccess:     config.MaxAccess,
		maxWorkers:    config.MaxWorkers,
		activeWorker:  make(chan struct{}, config.MaxWorkers),
		doneWorker:    make(chan struct{}, config.MaxAccess),
	}
	results := manager.Start()

	fmt.Println("")
	fmt.Printf("Total Access Count: %d\n", manager.maxAccess)
	fmt.Printf("Concurrency: %d\n", manager.maxWorkers)
	fmt.Printf("Total Time: %d msec\n", manager.totalElapsedMsec)

	for _, result := range results {
		fmt.Println("")
		fmt.Printf("[%s]\n", result.url)
		fmt.Printf("Success: %d\n", result.success)
		fmt.Printf("Failure: %d\n", result.failure)
		fmt.Printf("Average Response Time: %d msec\n", result.averageElapsedMsec())
		fmt.Printf("Minimum Response Time: %d msec\n", result.minimumElapsedMsec)
		fmt.Printf("Maximum Response Time: %d msec\n", result.maximumElapsedMsec)
	}
}
