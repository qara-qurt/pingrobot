package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pingrobot/models"
	"time"
)

const (
	WOKER_COUNT = 3
	INTERVAL    = time.Second * 5
	TG_TOKEN    = "6040490988:AAFdeU7w_369vpt7iuwyfN4tqhK8wnUcHlM"
	TG_CHAT_ID  = "958729521"
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
		if r.Error != nil {
			SendMessage(r.Info())
		}
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

func getUrl() string {
	return fmt.Sprintf("https://api.telegram.org/bot%s", TG_TOKEN)
}

func SendMessage(text string) error {
	var err error
	var response *http.Response

	url := fmt.Sprintf("%s/sendMessage", getUrl())

	body, _ := json.Marshal(map[string]string{
		"chat_id": TG_CHAT_ID,
		"text":    text,
	})

	response, err = http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	defer response.Body.Close()

	body, err = io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	return nil
}
