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

	adminSheets, err := services.NewAdminSheets(ctx, gc, cfg.AdminSpreadsheetId, cfg.AdminReadRange)
	if err != nil {
		log.Fatalf("fatal error %v", err)
	}

	secrets, err := adminSheets.GetSecrets(ctx)
	if err != nil {
		log.Fatalf("fatal error %v", err)
	}

	fmt.Println(secrets)

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
	wifiType := ""

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
						msg.Text = "printout"
					case "events":
						msg.Text = "events"
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
				if wifiType == "" {
					if update.Message.Text == "гостевой" {
						wifiType = "гостевой"
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "сеть Lan_Guest пароль lan123456")
						bot.Send(msg)

						wifiType = ""
					}
	
					if update.Message.Text == "коворкинг" {
						wifiType = "коворкинг"
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите номер, полученный от администратора")
						bot.Send(msg)
					}
				} else {
					if wifiType == "коворкинг" {
						var isValidCode bool

						for _, s := range secrets {
							decoded, err := encoder.Decode(s)
							if err != nil {
								msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Произошла ошибка генерации кода. Обратитесь к администратору")
								bot.Send(msg)
							}

							if update.Message.Text == decoded {
								isValidCode = true
								msg := tgbotapi.NewMessage(update.Message.Chat.ID, "сеть LAN пароль @lan2023")
								bot.Send(msg)

								wifiType = ""
							}
						}

						if !isValidCode {
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Пароль не верный, уточните у администратора")
							bot.Send(msg)
						}
					}
				}
			}
		}
	}
}