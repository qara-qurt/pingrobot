package main

import (
	"fmt"
	"log"
	"pingrobot/models"
	"pingrobot/storage"
	"pingrobot/telegram"
	"time"
)

const (
	WOKER_COUNT = 100
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
	stor, err := storage.InitDatabase()
	if err != nil {
		log.Fatal(err)
	}

	bot, err := telegram.InitTelegramBot(TG_TOKEN)
	if err != nil {
		log.Fatal(err)
	}

	user := make(chan models.UserInfo)
	go func() {
		if err := bot.CheckUpdates(user, stor); err != nil {
			log.Fatal(err)
		}
	}()

	for u := range user {
		checking(u, *bot, *stor)
	}

}

func checking(data models.UserInfo, bot telegram.TelegramBot, storage storage.Storage) {
	jobs := make(chan models.URL, len(urls))
	result := make(chan models.Result, len(urls))
	go createPool(jobs, result)
	go writeJobs(jobs, data)

	for r := range result {
		info := r.Info()
		fmt.Println(info)
		if r.Error != nil {
			bot.SendMessage(r.ChatId, info)
		}
	}
}

func writeJobs(jobs chan<- models.URL, data models.UserInfo) {
	for {
		for _, url := range data.URLs {
			urlR := models.NewURL(url, data.UserId, data.ChatId)
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
