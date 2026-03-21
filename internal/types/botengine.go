package types

type Services struct {
	CoworkersSheets CoworkersSheetsService
	Guests          GuestSheetsService
	BotLogs         BotLogsSheetsService
	Events          EventsService
	Subscriptions   SubscriptionsService

	CoworkingRegistrations CoworkingRegistrationsService
}
