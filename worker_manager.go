package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type WorkerManager struct {
	workers          []*Worker
	urls             []string
	basicAuthUser    string
	basicAuthPass    string
	maxAccess        int
	maxWorkers       int
	totalElapsedMsec time.Duration
	activeWorker     chan struct{}
	doneWorker       chan struct{}
}

func (manager *WorkerManager) Start() {
	watcher := &WorkerWatcher{doneWorker: manager.doneWorker}
	watcher.Start()
	manager.RunWorkers()
	watcher.WaitForFinish()
	manager.totalElapsedMsec = watcher.elapsedMsec
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

func (wm *WorkerManager) Cleanup() {
	close(wm.activeWorker)
	close(wm.doneWorker)
}

func (wm *WorkerManager) ShowResults() {

	fmt.Println("")
	fmt.Printf("Total Access Count: %d\n", wm.maxAccess)
	fmt.Printf("Concurrency: %d\n", wm.maxWorkers)
	fmt.Printf("Total Time: %d msec\n", wm.totalElapsedMsec)

	results := wm.CreateBaseResults()
	wm.CountStatusCode(results)
	wm.CalcElapsedTime(results)

	results.ShowResultsOfAll()
	results.ShowResultsOfEachURL()
}

func (wm *WorkerManager) CreateBaseResults() Results {
	results := Results{}
	for _, worker := range wm.workers {
		url := worker.request.URL.String()
		if _, exists := results[url]; exists == false {
			results[url] = &Result{}
		}
	}
	return results
}

func (wm *WorkerManager) CalcElapsedTime(results Results) {
	for _, worker := range wm.workers {
		url := worker.request.URL.String()
		results[url].sumElapsedMsec += worker.elapsedMsec

		if results[url].maximumElapsedMsec < worker.elapsedMsec {
			results[url].maximumElapsedMsec = worker.elapsedMsec
		}

		if results[url].minimumElapsedMsec > worker.elapsedMsec || results[url].minimumElapsedMsec == 0 {
			results[url].minimumElapsedMsec = worker.elapsedMsec
		}
	}
}

func (wm *WorkerManager) CountStatusCode(results Results) {
	for _, worker := range wm.workers {
		url := worker.request.URL.String()
		if strconv.Itoa(worker.statusCode)[0:1] == "2" {
			results[url].success++
		} else {
			results[url].failure++
		}
	}
}
