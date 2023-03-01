package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

var chatID int64

type TelegramBot struct {
	bot *tgbotapi.BotAPI
}

func InitTelegramBot(token string) (*TelegramBot, error) {
	tg, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return &TelegramBot{}, err
	}

	bot := &TelegramBot{
		bot: tg,
	}

	bot.bot.Debug = true

	return bot, nil
}

func (t TelegramBot) SendMessage(text string) {

	msg := tgbotapi.NewMessage(chatID, text)

	if _, err := t.bot.Send(msg); err != nil {
		log.Panic(err)
	}
}

func (t TelegramBot) CheckUpdates() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := t.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message.Command() == "start" {
			//save chatID
			chatID = update.Message.Chat.ID
			return
		}
	}

}
