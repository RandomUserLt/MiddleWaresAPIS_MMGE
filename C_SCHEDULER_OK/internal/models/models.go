package models

type Agenda struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Aligné sur ton Timetable
type Event struct {
	ID          string   `json:"id"`
	AgendaIDs   []string `json:"agendaIds"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Start       string   `json:"start"`       // RFC3339
	End         string   `json:"end"`         // RFC3339
	Location    string   `json:"location"`
	LastUpdate  string   `json:"lastUpdate"`  // RFC3339
}

