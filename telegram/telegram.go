package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"pingrobot/lib"
	"pingrobot/models"
	"pingrobot/storage"
	"strings"
)

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

func (t TelegramBot) SendMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)

	if _, err := t.bot.Send(msg); err != nil {
		return err
	}
	return nil
}

func (t TelegramBot) CheckUpdates(res chan<- models.UserInfo, stor *storage.Storage) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := t.bot.GetUpdatesChan(u)

	for update := range updates {
		switch update.Message.Command() {
		case "start":
			text := "Привет, я pingrobot\n" +
				"Что я умею - сохранять 'url'-ы вашего сервера и в случае если ваш сервер упадат я сообщу вам!\n" +
				"Как добавить urls - отправьте мне  /add -ваш url\n" +
				"Можете добавить несколько свойх серверов"

			if err := stor.Create(update.Message.From.ID); err != nil {
				return err
			}

			if err := t.SendMessage(update.Message.Chat.ID, text); err != nil {
				return err
			}
		case "add":
			userId := update.Message.From.ID
			chatId := update.Message.Chat.ID
			url := strings.Split(update.Message.Text, " ")[1]

			msg := "Отлично, теперь в трэкаю ваш сервер!"
			if ok := lib.ValidateURL(url); !ok {
				msg = "Введите коректный url аддрес"
			}

			if err := stor.Add(url, userId); err != nil {
				return err
			}

			if err := t.SendMessage(chatId, msg); err != nil {
				return err
			}

			user, err := stor.Get(userId)
			if err != nil {
				return err
			}

			res <- models.UserInfo{
				UserId: user.UserId,
				ChatId: chatId,
				URLs:   user.URLs,
			}
		}

	}
	return nil
}
