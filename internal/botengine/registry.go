package botengine

import (
	"strconv"
	"strings"

	"github.com/mkokoulin/LAN-coworking-bot/internal/state"
	"github.com/mkokoulin/LAN-coworking-bot/internal/types"
)


type UserRef struct {
	ID        int64
	Username  string // without "@", empty if hidden
	FirstName string
	LastName  string
	// можно потом добавить LanguageCode и т.п.
}

func (u UserRef) HasUsername() bool { return u.Username != "" }

func (u UserRef) AtUsername() string {
	if u.Username == "" {
		return ""
	}
	return "@" + u.Username
}

func (u UserRef) FullName() string {
	switch {
	case u.FirstName != "" && u.LastName != "":
		return u.FirstName + " " + u.LastName
	case u.FirstName != "":
		return u.FirstName
	case u.LastName != "":
		return u.LastName
	default:
		return ""
	}
}

func (u UserRef) Label() string {
	if u.Username != "" {
		return "@" + u.Username
	}
	if name := u.FullName(); name != "" {
		return name + " (id:" + strconv.FormatInt(u.ID, 10) + ")"
	}
	return "id:" + strconv.FormatInt(u.ID, 10)
}

type EventKind int

type Event struct {
	Kind            EventKind
	Command         string
	Text            string
	CallbackData    string
	CallbackQueryID string
	ChatID          int64
	MessageID       int
	InlineMessageID string
	FromUserName    string
	FromUserID      int64

	From UserRef
}

type FlowEntry struct {
	Flow types.Flow
	Step types.Step
}

type Registry struct {
	flows    map[types.Flow]map[types.Step]StepHandler
	commands map[string]FlowEntry
	cbPref   map[string]FlowEntry
	Store    state.Manager
}

const (
	EventCommand EventKind = iota
	EventText
	EventCallback
	MyChatMember
)

func NewRegistry(store state.Manager) *Registry {
	if store == nil {
		store = state.NewMemoryManager()
	}
	return &Registry{
		flows:    map[types.Flow]map[types.Step]StepHandler{},
		commands: map[string]FlowEntry{},
		cbPref:   map[string]FlowEntry{},
		Store:    store,
	}
}

func (r *Registry) RegisterFlow(flow types.Flow, steps map[types.Step]StepHandler) { r.flows[flow] = steps }
func (r *Registry) RegisterCommand(cmd string, entry FlowEntry)                     { r.commands[cmd] = entry }
func (r *Registry) RegisterCallbackPrefix(prefix string, entry FlowEntry)           { r.cbPref[prefix] = entry }

func (r *Registry) ResolveEntry(ev Event) (FlowEntry, bool) {
	switch ev.Kind {
	case EventCommand:
		if e, ok := r.commands[ev.Command]; ok {
			return e, true
		}
	case EventCallback:
		for pref, e := range r.cbPref {
			if strings.HasPrefix(ev.CallbackData, pref) {
				return e, true
			}
		}
	}
	return FlowEntry{}, false
}

func (r *Registry) Has(flow types.Flow, step types.Step) bool {
	if m, ok := r.flows[flow]; ok {
		_, ok2 := m[step]
		return ok2
	}
	return false
}
