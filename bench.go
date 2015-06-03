package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"net/http"
	"os"
	"time"
)

type WorkerManager struct {
	workers       []*Worker
	url           string
	basicAuthUser string
	basicAuthPass string
	activeWorker  chan struct{}
	doneWorker    chan struct{}
}

func (wm *WorkerManager) CreateWorker() (*Worker, error) {
	wm.activeWorker <- struct{}{}
	client := &http.Client{Timeout: time.Duration(100) * time.Second}
	request, err := http.NewRequest("GET", wm.url, nil)
	worker := &Worker{
		client:       *client,
		request:      *request,
		activeWorker: wm.activeWorker,
		doneWorker:   wm.doneWorker,
	}

	if wm.NeedToBasicAuthSet() {
		worker.SetBasicAuth(wm.basicAuthUser, wm.basicAuthPass)
	}

	wm.workers = append(wm.workers, worker)

	return worker, err
}

func (wm *WorkerManager) NeedToBasicAuthSet() bool {
	if wm.basicAuthUser != "" && wm.basicAuthPass != "" {
		return true
	}
	return false
}

func (wm *WorkerManager) WaitForWorkersToFinish() {
	for i := 0; i < cap(wm.doneWorker); i++ {
		<-wm.doneWorker
	}
}

func (wm *WorkerManager) Finish() {
	close(wm.activeWorker)
	close(wm.doneWorker)
}

type Worker struct {
	client       http.Client
	request      http.Request
	elapsedMsec  time.Duration
	activeWorker chan struct{}
	doneWorker   chan struct{}
}

func (w *Worker) SetBasicAuth(user string, password string) {
	w.request.SetBasicAuth(user, password)
}

func (w *Worker) Run() {
	defer func() {
		<-w.activeWorker
		w.doneWorker <- struct{}{}
	}()

	start := time.Now()
	response, err := w.client.Do(&w.request)
	w.elapsedMsec = time.Now().Sub(start) / time.Millisecond
	if err != nil {
		fmt.Printf("%sへのアクセスに失敗しました %s\n", w.request.URL, err)
		return
	}

	fmt.Printf("Response Time: %d msec, Status: %s\n", w.elapsedMsec, response.Status)
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
