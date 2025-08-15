package flows

import (
    "github.com/mkokoulin/LAN-coworking-bot/internal/botengine"
    flowsAbout "github.com/mkokoulin/LAN-coworking-bot/internal/flows/about"
    flowsBooking "github.com/mkokoulin/LAN-coworking-bot/internal/flows/booking"
    flowsEvents "github.com/mkokoulin/LAN-coworking-bot/internal/flows/events"
    flowsLanguage "github.com/mkokoulin/LAN-coworking-bot/internal/flows/language"
    flowsMenu "github.com/mkokoulin/LAN-coworking-bot/internal/flows/menu"
    flowsMeeting "github.com/mkokoulin/LAN-coworking-bot/internal/flows/meetingroom"
    flowsPrintout "github.com/mkokoulin/LAN-coworking-bot/internal/flows/printout"
    flowsStart "github.com/mkokoulin/LAN-coworking-bot/internal/flows/start"
    flowsWifi "github.com/mkokoulin/LAN-coworking-bot/internal/flows/wifi"
    flowsBar "github.com/mkokoulin/LAN-coworking-bot/internal/flows/bar"
    flowsKotolog "github.com/mkokoulin/LAN-coworking-bot/internal/flows/kotolog"
)

func RegisterAll(reg *botengine.Registry) {
    flowsStart.Register(reg)
    flowsLanguage.Register(reg)
    flowsWifi.Register(reg)
    flowsBooking.Register(reg)
    flowsMeeting.Register(reg)
    flowsPrintout.Register(reg)
    flowsEvents.Register(reg)
    flowsAbout.Register(reg)
    flowsMenu.Register(reg)
    flowsBar.Register(reg)
    flowsKotolog.Register(reg)
}
