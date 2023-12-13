package commands

import (
	"context"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mkokoulin/LAN-coworking-bot/internal/config"
)

func Wifi(ctx context.Context, update tgbotapi.Update, bot *tgbotapi.BotAPI, cfg *config.Config, args CommandsHandlerArgs) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	if !args.ChatState.IsWifiConfirmationProcess {
		if args.ChatState.Language == Languages[0].Lang {
			msg.Text = "Select the network options below: guest / coworking"

			msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("guest"),
					tgbotapi.NewKeyboardButton("coworking"),
				),
			)
		} else if args.ChatState.Language == Languages[1].Lang {
			msg.Text = "Выберите ниже варианты сети: гостевой / коворкинг"
	
			msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("гостевой"),
					tgbotapi.NewKeyboardButton("коворкинг"),
				),
			)
		}

		args.ChatState.IsWifiConfirmationProcess = true
		
		_, err := bot.Send(msg)
		return err
	}
		

	if !args.ChatState.IsAwaitingConfirmation {
		if update.Message.Text == "guest" || update.Message.Text == "гостевой" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			
			if args.ChatState.Language == Languages[0].Lang {
				msg.Text = fmt.Sprintf("Lan_Guest network password %s", cfg.GuestWifiPassword)
			} else if args.ChatState.Language == Languages[1].Lang {
				msg.Text = fmt.Sprintf("сеть Lan_Guest пароль %s", cfg.GuestWifiPassword)
			}

			args.ChatState.CurrentCommand = ""

			args.ChatState.IsWifiConfirmationProcess = false

			msg.ReplyMarkup = tgbotapi.ReplyKeyboardRemove{
				RemoveKeyboard: true,
				Selective: false,
			}

			_, err := bot.Send(msg)
			return err
		}

		if update.Message.Text == "coworking" || update.Message.Text == "коворкинг" {
			if args.ChatState.IsAuthorized {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

				if args.ChatState.Language == Languages[0].Lang {
					msg.Text = fmt.Sprintf("LAN network password %s", cfg.CoworkingWifiPassword)
				} else if args.ChatState.Language == Languages[1].Lang {
					msg.Text = fmt.Sprintf("сеть LAN пароль %s", cfg.CoworkingWifiPassword)
				}

				args.ChatState.CurrentCommand = ""

				args.ChatState.IsWifiConfirmationProcess = false

				msg.ReplyMarkup = tgbotapi.ReplyKeyboardRemove{
					RemoveKeyboard: true,
					Selective: false,
				}
				
				_, err := bot.Send(msg)
				return err
			}
			
				args.ChatState.IsAwaitingConfirmation = true

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

				if args.ChatState.Language == Languages[0].Lang {
					msg.Text = "Enter the number you received from the administrator"
				} else if args.ChatState.Language == Languages[1].Lang {
					msg.Text = "Введите номер, полученный от администратора"
				}
				
				_, err := bot.Send(msg)
				return err
			}
		} else {
			secrets, err := args.CoworkersSheets.GetUnusedSecrets(ctx)
			if err != nil {
				log.Fatalf("fatal error %v", err)
			}

			for _, s := range secrets {
				if update.Message.Text == s {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
					
					if args.ChatState.Language == Languages[0].Lang {
						msg.Text = fmt.Sprintf("LAN network password %s", cfg.CoworkingWifiPassword)
					} else if args.ChatState.Language == Languages[1].Lang {
						msg.Text = fmt.Sprintf("сеть LAN пароль %s", cfg.CoworkingWifiPassword)
					}

					args.ChatState.CurrentCommand = ""
					
					args.ChatState.IsWifiConfirmationProcess = false
					
					msg.ReplyMarkup = tgbotapi.ReplyKeyboardRemove{
						RemoveKeyboard: true,
						Selective: false,
					}

					args.ChatState.IsAwaitingConfirmation = false
					args.ChatState.IsAuthorized = true

					_, err := bot.Send(msg)
					return err
				}
			}

			if args.ChatState.IsAwaitingConfirmation {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

				if args.ChatState.Language == Languages[0].Lang {
					msg.Text = "The password is incorrect, check with the administrator"
				} else if args.ChatState.Language == Languages[1].Lang {
					msg.Text = "Пароль неверный, уточните у администратора"
				}
				
				_, err := bot.Send(msg)
				return err
			}
		}
	return nil
}