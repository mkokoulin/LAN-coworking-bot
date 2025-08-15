package types

// type EventKind int
// const (
//     EventCommand EventKind = iota
//     EventText
//     EventCallback
// )

type Event struct {
	ID          string `json:"id"`           // уникальный идентификатор (для URL/регистрации)
	Date        string `json:"date"`         // "2006-01-02" или "02.01.2006" — ты уже парсишь оба формата
	Name        string `json:"name"`         // краткое название события (plain text)
	Description string `json:"description"`  // описание (может быть с HTML)
	ShowForm    bool   `json:"showForm"`     // показывать в списке (фильтруешь по нему)
}
