package main

import (
	"net/http"
	"time"
)

type Worker struct {
	client       http.Client
	request      http.Request
	statusCode   int
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
	resp, _ := w.client.Do(&w.request)
	w.elapsedMsec = time.Now().Sub(start) / time.Millisecond
	defer resp.Body.Close()

	w.statusCode = resp.StatusCode
}
