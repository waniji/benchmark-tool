package main

import "flag"
import "fmt"
import "net/http"
import "os"
import "time"

type Worker struct {
	client  http.Client
	request http.Request
}

func CreateWorker(url string) (*Worker, error) {
	client := &http.Client{Timeout: time.Duration(100) * time.Second}
	request, err := http.NewRequest("GET", url, nil)
	worker := &Worker{client: *client, request: *request}

	return worker, err
}

func (w *Worker) SetBasicAuth(user string, password string) {
	w.request.SetBasicAuth(user, password)
}

func (w *Worker) Request() (*http.Response, time.Duration, error) {
	start := time.Now()
	response, err := w.client.Do(&w.request)
	elapsed_msec := time.Now().Sub(start) / time.Millisecond

	return response, elapsed_msec, err
}

func main() {
	var (
		url             string
		basic_auth_user string
		basic_auth_pass string
		max_count       int
		max_workers     int
	)

	flag.StringVar(&url, "url", "", "アクセスするURL")
	flag.IntVar(&max_count, "count", 1, "URLにアクセスする回数")
	flag.IntVar(&max_workers, "worker", 1, "同時アクセス数")
	flag.StringVar(&basic_auth_user, "basic-auth-user", "", "BASIC認証に使用するユーザー")
	flag.StringVar(&basic_auth_pass, "basic-auth-pass", "", "BASIC認証に使用するパスワード")
	flag.Parse()
	if url == "" {
		fmt.Printf("urlが不正です\n")
		os.Exit(1)
	}

	fmt.Printf("URL: %s\n", url)
	fmt.Printf("Total Access Count: %d\n", max_count)
	fmt.Printf("Concurrency: %d\n", max_workers)
	fmt.Println("--------------------------------------------------")

	worker := make(chan struct{}, max_workers)
	done := make(chan struct{}, max_count)
	ch := make(chan time.Duration, max_count)

	main_start := time.Now()
	for count := 0; count < max_count; count++ {

		worker <- struct{}{}

		go func() {
			defer func() {
				<-worker
				done <- struct{}{}
			}()

			worker, err := CreateWorker(url)
			if err != nil {
				fmt.Println(err)
				return
			}

			if basic_auth_user != "" && basic_auth_pass != "" {
				worker.SetBasicAuth(basic_auth_user, basic_auth_pass)
			}

			response, elapsed_msec, err := worker.Request()
			if err != nil {
				fmt.Printf("%sへのアクセスに失敗しました %s\n", url, err)
				return
			}

			ch <- elapsed_msec

			fmt.Printf("Response Time: %d msec, Status: %s\n", elapsed_msec, response.Status)
		}()
	}

	for i := 0; i < max_count; i++ {
		<-done
	}
	main_elapsed := time.Now().Sub(main_start) / time.Millisecond
	close(worker)
	close(done)
	close(ch)

	var total_elapsed time.Duration = 0
	var minimum_elapsed time.Duration = 0
	var maximum_elapsed time.Duration = 0
	for elapsed := range ch {
		total_elapsed += elapsed
		if maximum_elapsed < elapsed {
			maximum_elapsed = elapsed
		}
		if minimum_elapsed > elapsed || minimum_elapsed == 0 {
			minimum_elapsed = elapsed
		}
	}

	fmt.Println("--------------------------------------------------")
	fmt.Printf("Total Time           : %d msec\n", main_elapsed)
	fmt.Printf("Average Response Time: %d msec\n", total_elapsed/time.Duration(max_count))
	fmt.Printf("Minimum Response Time: %d msec\n", minimum_elapsed)
	fmt.Printf("Maximum Response Time: %d msec\n", maximum_elapsed)
}
