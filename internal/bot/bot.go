package bot

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/FatimaAlKaus/nparser"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

var (
	ErrIncorrectMangaId    = errors.New("Невалидный id манги")
	ErrMangaNotFound       = errors.New("Манга не найдена")
	ErrInternalServerError = errors.New("Произошла ошибка на сервере")
)

type Bot struct {
	bot     *tg.BotAPI
	nclient *nparser.Client
	logger  *logrus.Logger
}

func New(key string, client *nparser.Client, logger *logrus.Logger) (*Bot, error) {
	bot, err := tg.NewBotAPI(key)
	if err != nil {
		return nil, err
	}

	return &Bot{
		bot:     bot,
		nclient: client,
		logger:  logger,
	}, nil
}

func (b *Bot) Run() {

	b.logger.Infof("Authorized on account %s", b.bot.Self.UserName)

	u := tg.NewUpdate(0)
	u.Timeout = 10

	updates := b.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			b.logger.Infof("[%s]: %s", update.Message.From.UserName, update.Message.Text)

			mangaId, err := strconv.Atoi(update.Message.Text)
			if err != nil {
				b.bot.Send(tg.NewMessage(update.Message.Chat.ID, ErrIncorrectMangaId.Error()))
				continue
			}

			info, err := b.nclient.ComicInfo(mangaId)
			if err != nil {
				if errors.Is(err, nparser.ErrNotFound) {
					b.bot.Send(tg.NewMessage(update.Message.Chat.ID, ErrMangaNotFound.Error()))
					continue
				}
				b.bot.Send(tg.NewMessage(update.Message.Chat.ID, ErrInternalServerError.Error()))
				b.logger.Error(err)
				continue
			}

			b.bot.Send(tg.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%s\n%d страниц", info.Title.English, info.NumPages)))
			urls, err := b.nclient.GetURLPages(mangaId)
			if err != nil {
				if errors.Is(err, nparser.ErrNotFound) {
					b.bot.Send(tg.NewMessage(update.Message.Chat.ID, ErrMangaNotFound.Error()))
					continue
				} else {
					b.bot.Send(tg.NewMessage(update.Message.Chat.ID, ErrInternalServerError.Error()))
					b.logger.Error(err)
					continue
				}
			}

			photos := make([]interface{}, 0)
			for _, url := range urls {
				if len(photos) == 10 {
					media := tg.MediaGroupConfig{
						ChatID: update.Message.Chat.ID,
						Media:  photos,
					}
					b.bot.SendMediaGroup(media)

					photos = make([]interface{}, 0)
				}

				photos = append(photos, tg.NewInputMediaPhoto(tg.FileURL(url)))
			}
			if len(photos) > 0 {
				media := tg.MediaGroupConfig{
					ChatID: update.Message.Chat.ID,
					Media:  photos,
				}
				b.bot.SendMediaGroup(media)
			}
			b.bot.Send(tg.NewMessage(update.Message.Chat.ID, "Конец"))
		}
	}
}
