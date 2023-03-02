package models

import (
	"fmt"
	"pingrobot/client"
	"time"
)

type Result struct {
	UserId     int64
	ChatId     int64
	URL        string
	StatusCode int
	Timeout    time.Duration
	Error      error
}

func (r *Result) Info() string {
	if r.Error != nil {
		return fmt.Sprintf("[ERROR] - [%s] - %s", r.URL, r.Error.Error())
	}
	return fmt.Sprintf("[SUCCESS] - [%s] - Status: %d, Response Time: %s", r.URL, r.StatusCode, r.Timeout.String())
}

type URL struct {
	UserId int64
	ChatId int64
	URL    string
}

func NewURL(url string, userId, chatId int64) URL {
	return URL{
		URL:    url,
		UserId: userId,
		ChatId: chatId,
	}
}

func (u URL) Process() Result {
	result := Result{URL: u.URL, ChatId: u.ChatId}
	client := client.NewClient()
	now := time.Now()

	res, err := client.Client.Get(u.URL)
	if err != nil {
		result.Error = err
		return result
	}

	result.StatusCode = res.StatusCode
	result.Timeout = time.Since(now)
	return result
}

type UserInfo struct {
	UserId int64
	ChatId int64
	URLs   []string
}
