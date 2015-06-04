package main

import (
	"time"
)

type WorkerWatcher struct {
	doneWorker  chan struct{}
	allDone     chan struct{}
	elapsedMsec time.Duration
}

func (w *WorkerWatcher) Start() {
	w.allDone = make(chan struct{})
	go w.watch()
}

func (w *WorkerWatcher) WaitForFinish() {
	<-w.allDone
}

func (w *WorkerWatcher) watch() {
	start := time.Now()
	for i := 0; i < cap(w.doneWorker); i++ {
		<-w.doneWorker
	}
	w.allDone <- struct{}{}
	w.elapsedMsec = time.Now().Sub(start) / time.Millisecond
}
