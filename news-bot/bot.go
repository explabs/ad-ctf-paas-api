package news_bot

import (
	"fmt"
	"github.com/explabs/ad-ctf-paas-api/config"
)

func SendNews(filenames []string) error {
	token := config.Conf.Telegram.BotToken
	chatId := config.Conf.Telegram.ChatID
	if token == "" || chatId == "" {
		return fmt.Errorf("telecgram credentials is empty")
	}
	bot := TelegramBot{
		TelegramBotToken: token,
		ChatID:           chatId,
		NewsFolder:       "",
		Round:            0,
	}

	for _, filename := range filenames {
		if filename != "" {
			if err := bot.LoadMessage(filename); err != nil {
				return err
			}
			if err := bot.SendMessage(); err != nil {
				return err
			}
		}
	}
	return nil
}
