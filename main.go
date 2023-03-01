package main

import (
	"fmt"
	"log"
	"pingrobot/models"
	"pingrobot/telegram"
	"time"
)

const (
	WOKER_COUNT = 3
	INTERVAL    = time.Second * 5
	TG_TOKEN    = "6040490988:AAFdeU7w_369vpt7iuwyfN4tqhK8wnUcHlM"
)

var urls = []string{
	"https://mebel-f0c71.web.aapp/",
	"https://www.youtube.com",
	"https://www.google.com",
	"https://github.com/qara-qurt",
}

func main() {

	bot, err := telegram.InitTelegramBot(TG_TOKEN)
	if err != nil {
		log.Fatal(err.Error())
	}

	//after command '/start' continue code
	bot.CheckUpdates()

	jobs := make(chan models.URL, len(urls))
	result := make(chan models.Result, len(urls))

	go createPool(jobs, result)
	go writeJobs(jobs)

	//get result response and send telegram, show in terminal
	func() {
		for r := range result {
			info := r.Info()
			fmt.Println(info)
			if r.Error != nil {
				bot.SendMessage(info)
			}
		}
	}()

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
