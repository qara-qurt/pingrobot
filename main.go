package main

import (
	"fmt"
	"pingrobot/models"
	"time"
)

const (
	WOKER_COUNT = 3
	INTERVAL    = time.Second * 5
)

var urls = []string{
	"https://mebel-f0c71.web.aapp/",
	"https://www.youtube.com",
	"https://www.google.com",
	"https://github.com/qara-qurt",
}

func main() {
	jobs := make(chan models.URL, len(urls))
	result := make(chan models.Result, len(urls))

	go createPool(jobs, result)
	go getResponseInfo(result)
	writeJobs(jobs)
}

func getResponseInfo(result <-chan models.Result) {
	for r := range result {
		fmt.Println(r.Info())
	}
}

func writeJobs(jobs chan<- models.URL) {
	for {
		for _, url := range urls {
			urlR := models.NewURL(url)
			jobs <- urlR
		}
		time.Sleep(INTERVAL)
	}
}

func createPool(jobs <-chan models.URL, result chan<- models.Result) {
	for i := 0; i < WOKER_COUNT; i++ {
		go worker(jobs, result)
	}
}

func worker(jobs <-chan models.URL, result chan<- models.Result) {
	for job := range jobs {
		result <- job.Process()
	}
}
