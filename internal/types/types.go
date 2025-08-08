package types

type ChatStorage struct {
	Language               string
	IsAuthorized           bool
	IsBookingProcess       bool
	IsWifiProcess          bool
	IsAwaitingConfirmation bool
	CurrentCommand         string
}