package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"net/http"
	"os"
	"time"
)

type WorkerManager struct {
	worker        []Worker
	url           string
	basicAuthUser string
	basicAuthPass string
}

func (wm *WorkerManager) CreateWorker() (*Worker, error) {
	client := &http.Client{Timeout: time.Duration(100) * time.Second}
	request, err := http.NewRequest("GET", wm.url, nil)
	worker := &Worker{client: *client, request: *request}

	if wm.NeedToBasicAuthSet() {
		worker.SetBasicAuth(wm.basicAuthUser, wm.basicAuthPass)
	}

	return worker, err
}

func (wm *WorkerManager) NeedToBasicAuthSet() bool {
	if wm.basicAuthUser != "" && wm.basicAuthPass != "" {
		return true
	}
	return false
}

type Worker struct {
	client  http.Client
	request http.Request
}

func (w *Worker) SetBasicAuth(user string, password string) {
	w.request.SetBasicAuth(user, password)
}

func (w *Worker) Request() (*http.Response, time.Duration, error) {
	start := time.Now()
	response, err := w.client.Do(&w.request)
	elapsedMsec := time.Now().Sub(start) / time.Millisecond

	return response, elapsedMsec, err
}

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

	workerCh := make(chan struct{}, maxWorkers)
	done := make(chan struct{}, maxAccess)
	ch := make(chan time.Duration, maxAccess)

	manager := &WorkerManager{
		url:           url,
		basicAuthUser: basicAuthUser,
		basicAuthPass: basicAuthPass,
	}
	mainStart := time.Now()
	for count := 0; count < maxAccess; count++ {

		workerCh <- struct{}{}
		worker, err := manager.CreateWorker()
		if err != nil {
			fmt.Println(err)
			return
		}

		go func() {
			defer func() {
				<-workerCh
				done <- struct{}{}
			}()

			response, elapsedMsec, err := worker.Request()
			if err != nil {
				fmt.Printf("%sへのアクセスに失敗しました %s\n", url, err)
				return
			}

			ch <- elapsedMsec

			fmt.Printf("Response Time: %d msec, Status: %s\n", elapsedMsec, response.Status)
		}()
	}

	for i := 0; i < maxAccess; i++ {
		<-done
	}
	mainElapsedMsec := time.Now().Sub(mainStart) / time.Millisecond
	close(workerCh)
	close(done)
	close(ch)

	var totalElapsedMsec time.Duration = 0
	var minElapsedMsec time.Duration = 0
	var maxElapsedMsec time.Duration = 0
	for elapsed := range ch {
		totalElapsedMsec += elapsed
		if maxElapsedMsec < elapsed {
			maxElapsedMsec = elapsed
		}
		if minElapsedMsec > elapsed || minElapsedMsec == 0 {
			minElapsedMsec = elapsed
		}
	}

	fmt.Println("--------------------------------------------------")
	fmt.Printf("Total Time           : %d msec\n", mainElapsedMsec)
	fmt.Printf("Average Response Time: %d msec\n", totalElapsedMsec/time.Duration(maxAccess))
	fmt.Printf("Minimum Response Time: %d msec\n", minElapsedMsec)
	fmt.Printf("Maximum Response Time: %d msec\n", maxElapsedMsec)
}
