package flow

import (
	"github.com/mkokoulin/LAN-coworking-bot/internal/botengine"
	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
)

const (
	FlowDonation   types.Flow = "donation"
	DonationHome   types.Step = "donation:home"
	DonationCard   types.Step = "donation:card"
	DonationCopied types.Step = "donation:copied"
	DonationDone   types.Step = "donation:done"
)

// Замените на реальный номер карты
const cardNumber = "0000 0000 0000 0000"

// -------- Регистрация --------

func Register(reg *botengine.Registry) {
	// регистрируем RU переводы (EN — дефолт из кода)

	reg.RegisterFlow(FlowDonation, map[types.Step]botengine.StepHandler{
		DonationHome:   donationHome,
		DonationCard:   donationCard,
		DonationCopied: donationCopied,
		DonationDone:   donationDone,
	})

	// Команды/алиасы
	reg.RegisterCommand("donation", botengine.FlowEntry{Flow: FlowDonation, Step: DonationHome})
	reg.RegisterCommand("support", botengine.FlowEntry{Flow: FlowDonation, Step: DonationHome})

	// Все donation:* попадают в DonationHome, а там уже роутинг по CallbackData
	reg.RegisterCallbackPrefix("donation:", botengine.FlowEntry{Flow: FlowDonation, Step: DonationHome})
}
