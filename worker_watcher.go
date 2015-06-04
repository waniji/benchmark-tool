package main

import (
	"github.com/cheggaaa/pb"
	"time"
)

type WorkerWatcher struct {
	bar         pb.ProgressBar
	doneWorker  chan struct{}
	allDone     chan struct{}
	elapsedMsec time.Duration
}

func (w *WorkerWatcher) Start() {
	w.allDone = make(chan struct{})
	w.createProgressBar()
	go w.watch()
}

func (w *WorkerWatcher) WaitForFinish() {
	<-w.allDone
}

func (w *WorkerWatcher) createProgressBar() {
	bar := pb.New(cap(w.doneWorker))
	bar.SetMaxWidth(100)
	w.bar = *bar
}

func (w *WorkerWatcher) watch() {
	w.bar.Start()
	start := time.Now()
	for i := 0; i < cap(w.doneWorker); i++ {
		<-w.doneWorker
		w.bar.Increment()
	}
	w.allDone <- struct{}{}
	w.elapsedMsec = time.Now().Sub(start) / time.Millisecond
	w.bar.Finish()
}
