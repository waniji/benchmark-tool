package main

import (
	"time"
)

type Results []*Result

type Result struct {
	url                string
	sumElapsedMsec     time.Duration
	minimumElapsedMsec time.Duration
	maximumElapsedMsec time.Duration
	success            int
	failure            int
}

func (r *Result) averageElapsedMsec() time.Duration {
	return r.sumElapsedMsec / time.Duration(r.success+r.failure)
}
