package main

import (
	"log"

	"github.com/caarlos0/env/v11"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/panospet/go-perplexity-bot/internal/perplexity"
)

type config struct {
	PerplexityAPIKey string `env:"PERPLEXITY_API_KEY,required"`
	TelegramBotToken string `env:"TELEGRAM_BOT_TOKEN,required"`
}

func main() {
	var cfg config
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("%+v", err)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramBotToken)
	if err != nil {
		log.Fatalf("error creating new telegram bot: %v", err)
	}

	perplexitySrv := perplexity.NewService(
		cfg.PerplexityAPIKey,
	)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	log.Printf("Listening to bot messages...")

	for update := range updates {
		if update.Message == nil || !update.Message.IsCommand() {
			continue
		}

		switch update.Message.Command() {
		case "perplexity", "perp", "ask":
			question := update.Message.CommandArguments()
			if len(question) == 0 {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Please provide a question."))
				continue
			}

			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Thinking..."))
			answer, err := perplexitySrv.Ask(question)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error contacting Perplexity: "+err.Error()))
				continue
			}
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, answer))
		}
	}
}
