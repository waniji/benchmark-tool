package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	"time"
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

	fmt.Printf("URL: %s\n", url)
	fmt.Printf("Total Access Count: %d\n", maxAccess)
	fmt.Printf("Concurrency: %d\n", maxWorkers)
	fmt.Println("--------------------------------------------------")

	manager := &WorkerManager{
		url:           url,
		basicAuthUser: basicAuthUser,
		basicAuthPass: basicAuthPass,
		activeWorker:  make(chan struct{}, maxWorkers),
		doneWorker:    make(chan struct{}, maxAccess),
	}

	mainStart := time.Now()
	for count := 0; count < maxAccess; count++ {

		worker, err := manager.CreateWorker()
		if err != nil {
			fmt.Println(err)
			return
		}

		go worker.Run()
	}

	manager.WaitForWorkersToFinish()
	manager.Finish()
	mainElapsedMsec := time.Now().Sub(mainStart) / time.Millisecond

	var totalElapsedMsec time.Duration = 0
	var minElapsedMsec time.Duration = 0
	var maxElapsedMsec time.Duration = 0

	for _, worker := range manager.workers {
		totalElapsedMsec += worker.elapsedMsec
		if maxElapsedMsec < worker.elapsedMsec {
			maxElapsedMsec = worker.elapsedMsec
		}
		if minElapsedMsec > worker.elapsedMsec || minElapsedMsec == 0 {
			minElapsedMsec = worker.elapsedMsec
		}
	}

	fmt.Println("--------------------------------------------------")
	fmt.Printf("Total Time           : %d msec\n", mainElapsedMsec)
	fmt.Printf("Average Response Time: %d msec\n", totalElapsedMsec/time.Duration(maxAccess))
	fmt.Printf("Minimum Response Time: %d msec\n", minElapsedMsec)
	fmt.Printf("Maximum Response Time: %d msec\n", maxElapsedMsec)
}
