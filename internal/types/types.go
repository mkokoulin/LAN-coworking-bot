package types

type ChatStorage struct {
	CurrentCommand string
	IsAwaitingConfirmation bool
	IsAuthorized bool
	Language string
	IsBookingProcess bool
	IsWifiProcess bool
}