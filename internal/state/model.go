package state

// Flow: имя команды (wifi, booking...)
// Step: подшаг в этом флоу
// Data: временные данные между шагами
type ChatStorage struct {
	ID        int64                  `bson:"_id"`
	Language  string                 `bson:"language"`
	Flow      string                 `bson:"flow"`
	Step      string                 `bson:"step"`
	Data      map[string]interface{} `bson:"data"`
	Attempts  int                    `bson:"attempts"`
	IsAuth    bool                   `bson:"is_authorized"`
}
