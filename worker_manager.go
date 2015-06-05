package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type WorkerManager struct {
	workers       []*Worker
	urls          []string
	basicAuthUser string
	basicAuthPass string
	maxAccess     int
	maxWorkers    int
	result        Result
	activeWorker  chan struct{}
	doneWorker    chan struct{}
}

func (manager *WorkerManager) Start() {
	watcher := &WorkerWatcher{doneWorker: manager.doneWorker}
	watcher.Start()
	manager.RunWorkers()
	watcher.WaitForFinish()

	manager.result.totalElapsedMsec = watcher.elapsedMsec
	manager.CountStatusCode()
	manager.CalcElapsedTime()
	manager.Cleanup()
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
	request, err := http.NewRequest("GET", wm.SelectUrl(), nil)
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

func (wm *WorkerManager) SelectUrl() string {
	if len(wm.urls) == 1 {
		return wm.urls[0]
	}
	url := wm.urls[0]
	wm.urls = wm.urls[1:]
	wm.urls = append(wm.urls, url)
	return url
}

func (wm *WorkerManager) NeedToBasicAuthSet() bool {
	if wm.basicAuthUser != "" && wm.basicAuthPass != "" {
		return true
	}
	return false
}

func (wm *WorkerManager) CalcElapsedTime() {

	var sumTotalElapsedMsec time.Duration
	for _, worker := range wm.workers {
		sumTotalElapsedMsec += worker.elapsedMsec
		if wm.result.maximumElapsedMsec < worker.elapsedMsec {
			wm.result.maximumElapsedMsec = worker.elapsedMsec
		}
		if wm.result.minimumElapsedMsec > worker.elapsedMsec || wm.result.minimumElapsedMsec == 0 {
			wm.result.minimumElapsedMsec = worker.elapsedMsec
		}
	}
	wm.result.averageElapsedMsec = sumTotalElapsedMsec / time.Duration(wm.maxAccess)
}

func (wm *WorkerManager) CountStatusCode() {
	for _, worker := range wm.workers {
		if strconv.Itoa(worker.statusCode)[0:1] == "2" {
			wm.result.success++
		} else {
			wm.result.failure++
		}
	}
}

func (wm *WorkerManager) Cleanup() {
	close(wm.activeWorker)
	close(wm.doneWorker)
}
