package news_bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

type TelegramBot struct {
	TelegramBotToken string
	ChatID           string
	NewsFolder       string
	Round            int
	Text             []byte
}

type TelegramMessage struct {
	ChatID              string `json:"chat_id"`
	Text                string `json:"text"`
	DisableNotification bool   `json:"disable_notification"`
	ParseMode           string `json:"parse_mode"`
}

func (bot *TelegramBot) SendMessage() error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", bot.TelegramBotToken)
	body := TelegramMessage{
		ChatID:              bot.ChatID,
		Text:                string(bot.Text),
		DisableNotification: true,
		ParseMode:           "MarkdownV2",
	}

	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	return nil
}

func (bot *TelegramBot) LoadMessage(fileName string) error {
	file, err := os.Open(filepath.Join(bot.NewsFolder, fileName))
	if err != nil {
		return err
	}

	bot.Text, err = ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	if err = file.Close(); err != nil {
		return err
	}
	return nil
}
