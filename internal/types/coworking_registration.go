package types

import "time"

type RegistrationStatus string

const (
	RegistrationPending   RegistrationStatus = "pending"
	RegistrationApproved  RegistrationStatus = "approved"
	RegistrationRejected  RegistrationStatus = "rejected"
	RegistrationClarify   RegistrationStatus = "clarification_required"
)

type RegistrationRequestType string

const (
	RequestTypeNew    RegistrationRequestType = "new"
	RequestTypeRelink RegistrationRequestType = "relink"
)

type CoworkingRegistration struct {
	ID               string                  `bson:"_id,omitempty"`
	ChatID           int64                   `bson:"chat_id"`
	TelegramUserID   int64                   `bson:"telegram_user_id"`
	TelegramUsername string                  `bson:"telegram_username,omitempty"`

	FullName   string `bson:"full_name"`
	Phone      string `bson:"phone,omitempty"`
	TariffCode string `bson:"tariff_code,omitempty"`

	RequestType  RegistrationRequestType `bson:"request_type"`
	Status       RegistrationStatus      `bson:"status"`
	AdminComment string                  `bson:"admin_comment,omitempty"`

	CreatedAt  time.Time  `bson:"created_at"`
	UpdatedAt  time.Time  `bson:"updated_at"`
	ApprovedAt *time.Time `bson:"approved_at,omitempty"`
	ApprovedBy int64      `bson:"approved_by,omitempty"`
}