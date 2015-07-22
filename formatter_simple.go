package main

import (
	"fmt"
)

type FormatterSimple struct{}

func (f *FormatterSimple) Print(results Results) {
	for _, result := range results {
		f.PrintResult(result)
	}
}

func (f *FormatterSimple) PrintResult(r Result) {
	fmt.Println("")
	fmt.Printf("[%s]\n", r.url)
	fmt.Printf("Success: %d\n", r.success)
	fmt.Printf("Failure: %d\n", r.failure)
	fmt.Printf("Average Response Time: %d msec\n", r.averageElapsedMsec())
	fmt.Printf("Minimum Response Time: %d msec\n", r.minimumElapsedMsec)
	fmt.Printf("Maximum Response Time: %d msec\n", r.maximumElapsedMsec)
}
