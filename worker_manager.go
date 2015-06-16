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

func (manager *WorkerManager) Start() Results {
	watcher := &WorkerWatcher{doneWorker: manager.doneWorker}
	watcher.Start()
	manager.RunWorkers()
	watcher.WaitForFinish()
	manager.totalElapsedMsec = watcher.elapsedMsec
	manager.Cleanup()

	return manager.AggregateWorkerResult()
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

func (wm *WorkerManager) AggregateWorkerResult() Results {
	resultMap := make(map[string]*Result)

	for _, worker := range wm.workers {
		url := worker.request.URL.String()
		if _, exists := resultMap[url]; exists == false {
			resultMap[url] = &Result{}
		}

		resultMap[url].sumElapsedMsec += worker.elapsedMsec

		if resultMap[url].maximumElapsedMsec < worker.elapsedMsec {
			resultMap[url].maximumElapsedMsec = worker.elapsedMsec
		}

		if resultMap[url].minimumElapsedMsec > worker.elapsedMsec || resultMap[url].minimumElapsedMsec == 0 {
			resultMap[url].minimumElapsedMsec = worker.elapsedMsec
		}

		if strconv.Itoa(worker.statusCode)[0:1] == "2" {
			resultMap[url].success++
		} else {
			resultMap[url].failure++
		}
	}

	allResult := &Result{url: "all"}
	for _, result := range resultMap {
		allResult.success += result.success
		allResult.failure += result.failure
		allResult.sumElapsedMsec += result.sumElapsedMsec
		if allResult.minimumElapsedMsec > result.minimumElapsedMsec || allResult.minimumElapsedMsec == 0 {
			allResult.minimumElapsedMsec = result.minimumElapsedMsec
		}
		if allResult.maximumElapsedMsec < result.maximumElapsedMsec || allResult.maximumElapsedMsec == 0 {
			allResult.maximumElapsedMsec = result.maximumElapsedMsec
		}
	}

	results := Results{}
	results = append(results, allResult)
	for key, value := range resultMap {
		value.url = key
		results = append(results, value)
	}

	return results
}
