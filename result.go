package main

import (
	"errors"
	"sort"
	"time"
)

type Results []Result

type Result struct {
	url                string
	sumElapsedMsec     time.Duration
	minimumElapsedMsec time.Duration
	maximumElapsedMsec time.Duration
	success            int
	failure            int
}

func (r Result) averageElapsedMsec() time.Duration {
	return r.sumElapsedMsec / time.Duration(r.success+r.failure)
}

type ResultsSorter struct {
	results Results
	less    func(r1, r2 Result) bool
}

func CreateSorter(sortKey string) (ResultsSorter, error) {
	rs := ResultsSorter{}

	switch sortKey {
	case "success":
		rs.less = func(r1, r2 Result) bool {
			return r1.success < r2.success
		}
	case "failure":
		rs.less = func(r1, r2 Result) bool {
			return r1.failure < r2.failure
		}
	case "maximum":
		rs.less = func(r1, r2 Result) bool {
			return r1.maximumElapsedMsec < r2.maximumElapsedMsec
		}
	case "minimum":
		rs.less = func(r1, r2 Result) bool {
			return r1.minimumElapsedMsec < r2.minimumElapsedMsec
		}
	case "average":
		rs.less = func(r1, r2 Result) bool {
			return r1.averageElapsedMsec() < r2.averageElapsedMsec()
		}
	default:
		return rs, errors.New("sortが不正です: " + sortKey)
	}

	return rs, nil
}

func (rs ResultsSorter) Sort(results Results) {
	rs.results = results
	sort.Sort(rs)
}

func (rs ResultsSorter) Len() int {
	return len(rs.results)
}

func (rs ResultsSorter) Swap(i, j int) {
	rs.results[i], rs.results[j] = rs.results[j], rs.results[i]
}

func (rs ResultsSorter) Less(i, j int) bool {
	return rs.less(rs.results[i], rs.results[j])
}
