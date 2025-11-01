package models

type Event struct {
	ID          string   `json:"id"`
	AgendaIDs   []string `json:"agendaIds"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Start       string   `json:"start"`      // RFC3339
	End         string   `json:"end"`        // RFC3339
	Location    string   `json:"location"`
	LastUpdate  string   `json:"lastUpdate"` // RFC3339
}

// Un changement élémentaire
type Change struct {
	Field   string `json:"field"`   // "location", "start", "end", "title", "description"
	Old     string `json:"old"`
	New     string `json:"new"`
	Kind    string `json:"kind"`    // "room_change", "time_change", "title_change", "description_change"
}

// Message d’alerte prêt à être envoyé au service d’emails
type AlertMessage struct {
	Type        string   `json:"type"`        // "event_changed" | "event_new"
	EventID     string   `json:"eventId"`
	AgendaIDs   []string `json:"agendaIds"`
	Title       string   `json:"title"`
	Start       string   `json:"start"`
	End         string   `json:"end"`
	Location    string   `json:"location"`
	Changes     []Change `json:"changes"`     // liste des changements
	EmailText   string   `json:"emailText"`   // rendu prêt-à-insérer dans un mail
}
