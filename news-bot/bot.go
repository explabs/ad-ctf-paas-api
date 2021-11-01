package news_bot

import (
	"os"
)

func SendNews(filenames []string) error {
	bot := TelegramBot{
		TelegramBotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		ChatID:           os.Getenv("CHAT_ID"),
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
