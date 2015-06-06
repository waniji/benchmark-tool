package main

import (
	"fmt"
	"time"
)

type Results map[string]*Result

func (results *Results) ShowResultsOfAll() {
	var allResult Result
	for _, result := range *results {
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
	fmt.Println("")
	fmt.Println("[all]")
	allResult.ShowStausCount()
	allResult.ShowElapsedTime()
}

func (results *Results) ShowResultsOfEachURL() {
	for url, result := range *results {
		fmt.Println("")
		fmt.Printf("[%s]\n", url)
		result.ShowStausCount()
		result.ShowElapsedTime()
	}
}

type Result struct {
	sumElapsedMsec     time.Duration
	minimumElapsedMsec time.Duration
	maximumElapsedMsec time.Duration
	success            int
	failure            int
}

func (r *Result) averageElapsedMsec() time.Duration {
	return r.sumElapsedMsec / time.Duration(r.success+r.failure)
}

func (r *Result) ShowStausCount() {
	fmt.Printf("Success: %d\n", r.success)
	fmt.Printf("Failure: %d\n", r.failure)
}

func (r *Result) ShowElapsedTime() {
	fmt.Printf("Average Response Time: %d msec\n", r.averageElapsedMsec())
	fmt.Printf("Minimum Response Time: %d msec\n", r.minimumElapsedMsec)
	fmt.Printf("Maximum Response Time: %d msec\n", r.maximumElapsedMsec)
}
