package types

// ADD: модели для истории/состояния (json-friendly)
type BarState struct {
	Cart      map[string]int `json:"cart,omitempty"`
	Buyer     string         `json:"buyer,omitempty"`
	Serve     string         `json:"serve,omitempty"` // "pickup" | "tozone"
	Zone      string         `json:"zone,omitempty"`  // "coworking" | "cafe" | "street"
	Currency  string         `json:"currency,omitempty"`
	Notes     string         `json:"notes,omitempty"`
	OrderID   string         `json:"order_id,omitempty"` // последний, если есть
	UpdatedAt int64          `json:"updated_at,omitempty"`
}

type CommandEntry struct {
	Cmd     string `json:"cmd"`
	AtUnix  int64  `json:"at"`
	Payload string `json:"payload,omitempty"`
}

type BarOrderItem struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Qty      int    `json:"qty"`
	PriceAMD int    `json:"price_amd"`
	SumAMD   int    `json:"sum_amd"`
}

type BarOrder struct {
	OrderID   string         `json:"order_id"`
	Buyer     string         `json:"buyer"`
	Items     []BarOrderItem `json:"items"`
	TotalAMD  int            `json:"total_amd"`
	Serve     string         `json:"serve"`
	Zone      string         `json:"zone"`
	Notes     string         `json:"notes,omitempty"`
	CreatedAt int64          `json:"created_at"`
}
