package main

import (
	"net/http"
	"time"
)

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
	w.client.Do(&w.request)
	w.elapsedMsec = time.Now().Sub(start) / time.Millisecond
}

type Result struct {
	totalElapsedMsec   time.Duration
	averageElapsedMsec time.Duration
	minimumElapsedMsec time.Duration
	maximumElapsedMsec time.Duration
}
