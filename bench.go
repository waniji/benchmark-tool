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
