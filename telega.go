package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type Config struct {
	Token string `yaml:"token"`
}

func loadConfig() Config {
	filename := "config.yml"
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Cannot open config file %s, reason: %s", filename, err)
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatalf("Cannot decode config file %s, reason: %s", filename, err)
	}
	return cfg
}

func main() {
	cfg := loadConfig()
	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		message, err := bot.Send(msg)
		if err != nil {
			log.Printf("Error on sending message %s to %s", msg.Text, err)
		} else {
			log.Printf("[%s -> %s] %s", message.From.UserName, message.Chat.UserName, msg.Text)
		}
	}
}
