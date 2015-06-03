package main

import (
	"fmt"
	"net/http"
	"time"
)

type WorkerManager struct {
	workers       []*Worker
	url           string
	basicAuthUser string
	basicAuthPass string
	maxAccess     int
	maxWorkers    int
	result        Result
	activeWorker  chan struct{}
	doneWorker    chan struct{}
}

func (wm *WorkerManager) Start() {
	start := time.Now()
	wm.RunWorkers()
	wm.WaitForWorkersToFinish()
	wm.result.totalElapsedMsec = time.Now().Sub(start) / time.Millisecond
	wm.CalcElapsedTime()
	wm.Cleanup()
}

func (wm *WorkerManager) RunWorkers() {
	for count := 0; count < wm.maxAccess; count++ {

		wm.activeWorker <- struct{}{}
		worker, err := wm.CreateWorker()
		if err != nil {
			fmt.Println(err)
			return
		}
		go worker.Run()
	}
}

func (wm *WorkerManager) CreateWorker() (*Worker, error) {
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

func (wm *WorkerManager) CalcElapsedTime() {
	var result Result
	for _, worker := range wm.workers {
		result.totalElapsedMsec += worker.elapsedMsec
		if result.maximumElapsedMsec < worker.elapsedMsec {
			result.maximumElapsedMsec = worker.elapsedMsec
		}
		if result.minimumElapsedMsec > worker.elapsedMsec || result.minimumElapsedMsec == 0 {
			result.minimumElapsedMsec = worker.elapsedMsec
		}
	}
	result.averageElapsedMsec = wm.result.totalElapsedMsec / time.Duration(wm.maxAccess)
	wm.result = result
}

func (wm *WorkerManager) Cleanup() {
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

type Result struct {
	totalElapsedMsec   time.Duration
	averageElapsedMsec time.Duration
	minimumElapsedMsec time.Duration
	maximumElapsedMsec time.Duration
}
