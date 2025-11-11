package models

// Event represents a normalized VEVENT extracted from an iCal feed.
type Event struct {
	ID          string   `json:"id"`          // iCal UID
	AgendaIDs   []string `json:"agendaIds"`   // resources used to fetch
	Title       string   `json:"title"`       // SUMMARY
	Description string   `json:"description"` // DESCRIPTION (raw text)
	Start       string   `json:"start"`       // RFC3339
	End         string   `json:"end"`         // RFC3339
	Location    string   `json:"location"`    // LOCATION
	LastUpdate  string   `json:"lastUpdate"`  // LAST-MODIFIED (RFC3339)
}
