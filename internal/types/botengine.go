package types

import "github.com/mkokoulin/LAN-coworking-bot/internal/services"

type Services struct {
	CoworkersSheets *services.CoworkersSheetService
	GuestSheets     *services.GuestsSheetService
	BotLogsSheets   *services.BotLogsSheetService
}
