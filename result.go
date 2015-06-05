package main

import (
	"fmt"
	"time"
)

type Result struct {
	totalElapsedMsec   time.Duration
	averageElapsedMsec time.Duration
	minimumElapsedMsec time.Duration
	maximumElapsedMsec time.Duration
	success            int
	failure            int
}

func (r *Result) ShowStausCount() {
	fmt.Println("")
	fmt.Printf("Success: %d\n", r.success)
	fmt.Printf("Failure: %d\n", r.failure)
}

func (r *Result) ShowElapsedTime() {
	fmt.Println("")
	fmt.Printf("Total Time           : %d msec\n", r.totalElapsedMsec)
	fmt.Printf("Average Response Time: %d msec\n", r.averageElapsedMsec)
	fmt.Printf("Minimum Response Time: %d msec\n", r.minimumElapsedMsec)
	fmt.Printf("Maximum Response Time: %d msec\n", r.maximumElapsedMsec)
}
