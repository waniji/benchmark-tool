package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
)

func bench(c *cli.Context) {

	var config Config
	var formatter Formatter
	var sorter ResultsSorter
	var err error

	if err = config.Create(c); err != nil {
		fmt.Println(err)
		cli.ShowAppHelp(c)
		os.Exit(1)
	}

	formatter, err = CreateFormatter(config.Format)
	if err != nil {
		fmt.Println(err)
		cli.ShowAppHelp(c)
		os.Exit(1)
	}

	sorter, err = CreateSorter(config.Sort)
	if err != nil {
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

	sorter.Sort(results[1:])

	fmt.Println("")
	fmt.Printf("Total Access Count: %d\n", manager.maxAccess)
	fmt.Printf("Concurrency: %d\n", manager.maxWorkers)
	fmt.Printf("Total Time: %d msec\n", manager.totalElapsedMsec)

	formatter.Print(results)
}
