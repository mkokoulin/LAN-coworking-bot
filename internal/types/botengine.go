package types

type Services struct {
	CoworkersSheets CoworkersSheetsService
	GuestSheets     GuestSheetsService
	BotLogsSheets   BotLogsSheetsService
	Events         	EventsService
	Subscriptions   SubscriptionsService
}
