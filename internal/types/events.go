package types

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type Event struct {
	ID          string `json:"id"`          
	Date        string `json:"date"`        
	Name        string `json:"name"`        
	Description string `json:"description"`
	ShowForm    bool   `json:"showForm"`
	Capacity    int    `json:"capacity"` 
	ExternalLink string `json:"externalLink"`
}

// UnmarshalJSON — ест capacity как число ИЛИ строку.
// "", null, отсутствие поля => 0 (без лимита).
func (e *Event) UnmarshalJSON(data []byte) error {
	var aux struct {
		ID          string          `json:"id"`
		Date        string          `json:"date"`
		Name        string          `json:"name"`
		Description string          `json:"description"`
		ShowForm    bool            `json:"showForm"`
		Capacity    json.RawMessage `json:"capacity"`
		ExternalLink string          `json:"externalLink"`
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	e.ID = aux.ID
	e.Date = aux.Date
	e.Name = aux.Name
	e.Description = aux.Description
	e.ShowForm = aux.ShowForm
	e.ExternalLink = aux.ExternalLink

	// capacity отсутствует / null => 0
	if len(aux.Capacity) == 0 || string(aux.Capacity) == "null" {
		e.Capacity = 0
		return nil
	}

	// строка
	if aux.Capacity[0] == '"' {
		var s string
		if err := json.Unmarshal(aux.Capacity, &s); err != nil {
			return fmt.Errorf("capacity string: %w", err)
		}
		s = strings.TrimSpace(s)
		if s == "" {
			e.Capacity = 0 // пустая строка = без лимита
			return nil
		}
		n, err := strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("capacity atoi: %w", err)
		}
		e.Capacity = n
		return nil
	}

	// число
	if err := json.Unmarshal(aux.Capacity, &e.Capacity); err == nil {
		return nil
	}
	// на крайний случай — float -> int
	var f float64
	if err := json.Unmarshal(aux.Capacity, &f); err == nil {
		e.Capacity = int(f)
		return nil
	}

	// если совсем что-то экзотическое — трактуем как без лимита
	e.Capacity = 0
	return nil
}
