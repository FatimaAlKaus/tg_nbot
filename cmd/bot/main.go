package main

import (
	"flag"
	"main/internal/bot"
	"main/internal/config"

	"github.com/FatimaAlKaus/nparser"
	"github.com/sirupsen/logrus"
)

var _configPath = flag.String("config", "config/dev.yml", "path to the config file")

func main() {
	logger := logrus.New()
	cfg, err := config.Load(*_configPath)
	if err != nil {
		logger.Fatal(err)
	}
	bot, err := bot.New(cfg.BotToken, &nparser.Client{}, logrus.New())
	if err != nil {
		logger.Fatalf("failed to inut bot: %v", err)
	}

	// TODO Graceful shutdown
	bot.Run()
}
