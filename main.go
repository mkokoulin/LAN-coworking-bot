package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mkokoulin/LAN-coworking-bot/internal/config"
	"github.com/mkokoulin/LAN-coworking-bot/internal/helpers/encoder"
	"github.com/mkokoulin/LAN-coworking-bot/internal/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	ctx := context.Background()

	cfg, err := config.New()
	if err != nil {
		log.Fatalln(err)
		return
	}

	gc, err := services.NewGoogleClient(ctx, cfg.GoogleCloudConfig, cfg.Scope)
	if err != nil {
		log.Fatalf("fatal error %v", err)
	}

	coworkersSheets, err := services.NewCoworkersSheets(ctx, gc, cfg.CoworkersSpreadsheetId, cfg.CoworkersReadRange)
	if err != nil {
		log.Fatalf("fatal error %v", err)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	currentCommand := ""
	isAwaitingConfirmation := false

	for update := range updates {
		if update.Message != nil {
			if update.Message.IsCommand() {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
				currentCommand = update.Message.Command()

				switch currentCommand {
					case "start":
						msg.Text = "start"
					case "wifi":
						msg.Text = "Выберите вариант сети: гостевой / коворкинг"
						msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
							tgbotapi.NewKeyboardButtonRow(
								tgbotapi.NewKeyboardButton("гостевой"),
								tgbotapi.NewKeyboardButton("коворкинг"),
							),
						)
					case "meetingrooom":
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
						msg.Text = "Напишите дату и интервал времени, на который вы хотите забронировать комнату для переговоров в формате yyyy-mm-dd hh:mm - hh:mm"
					case "printout":
						msg.Text = "Отправьте документы для распечатки в аккаунт @lan_yerevan (администратору) и уточните у него стоимость услуги"
					case "events":
						msg.ParseMode = "html"
						msg.Text = "У нас проходит большое количество разнообразных мероприятий, анонсы событий мы публикуем в наших социальных сетях: <a href='https://www.instagram.com/lan_yerevan/'>Instagram</a> и <a href='https://t.me/lan_yerevan'>Telegram</a>. Подписывайтесь, чтобы быть в курсе классных событий 🎉. Актуальный список мероприятий и бронирование ведется через <a href='https://taplink.cc/lan_yerevan'>taplink</a>"
					case "about":
						msg.Text = "about"
					default:
						msg.Text = "I don't know that command"
				}
	
				if _, err := bot.Send(msg); err != nil {
					log.Print(err)
				}

				continue
			}

			if currentCommand == "meetingrooom" {
				if update.Message.Text == "" {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
					msg.Text = "Сообщение не может быть пустым"
					bot.Send(msg)

					continue
				}

				msgTo := tgbotapi.NewMessage(5701365900, fmt.Sprintf("Пользователь @%s просит забронировать переговорку - %s", update.Message.Chat.UserName, update.Message.Text))

				bot.Send(msgTo)
			}

			if currentCommand == "wifi" {
				if !isAwaitingConfirmation {
					if update.Message.Text == "гостевой" {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "сеть Lan_Guest пароль lan123456")
						bot.Send(msg)
					}
	
					if update.Message.Text == "коворкинг" {
						coworker, err := coworkersSheets.GetCoworker(ctx, fmt.Sprintf("@%s", update.Message.Chat.UserName))
						if err != nil {
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Произошла ошибка генерации кода. Обратитесь к администратору")
							bot.Send(msg)
						}

						if coworker.Telegram != "" {
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, "сеть LAN пароль @lan2023")
							bot.Send(msg)
							continue
						}
						
						isAwaitingConfirmation = true
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите номер, полученный от администратора")
						bot.Send(msg)
					}
				} else {
						unusedSecrets, err := coworkersSheets.GetUnusedSecrets(ctx)
						if err != nil {
							log.Fatalf("fatal error %v", err)
						}

						for _, s := range unusedSecrets {
							decoded, err := encoder.Decode(s)
							if err != nil {
								msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Произошла ошибка генерации кода. Обратитесь к администратору")
								bot.Send(msg)
							}

							if update.Message.Text == decoded {
								msg := tgbotapi.NewMessage(update.Message.Chat.ID, "сеть LAN пароль @lan2023")
								bot.Send(msg)

								newCoworker := services.Coworker{
									Secret: s,
									Telegram: fmt.Sprintf("@%s", update.Message.Chat.UserName),
								}
								err := coworkersSheets.UpdateCoworker(ctx, newCoworker)
								if err != nil {
									msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Произошла ошибка генерации кода. Обратитесь к администратору")
									bot.Send(msg)
								}
								isAwaitingConfirmation = false

								continue
							}

						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Пароль не верный, уточните у администратора")
						bot.Send(msg)
					}
				}
			}
		}
	}
}