package types

import "github.com/google/uuid"

type ChatStorage struct {
	ID             uuid.UUID `json:"id"`
	ChatID int64 `json:"chat_id"`
	CurrentCommand string `json:"current_command"`
	PreviousCommand string `json:"previous_command"`
	IsAwaitingConfirmation bool `json:"is_awaiting_confirmation"`
	IsAuthorized bool `json:"is_authorized"`
	Language string `json:"language"`
	IsLanguageSelectionProcess bool `json:"is_language_selection_process"`
	IsMeetingroomBookingProcess bool `json:"is_meetingroom_booking_process"`
	IsWifiConfirmationProcess bool `json:"is_wifi_confirmation_process"`
}