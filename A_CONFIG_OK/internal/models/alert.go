package models

// Alert represents a notification rule bound to an agenda and a recipient.
// Condition examples (per TP): "always", "room_change", "time_change", "teacher_change".
type Alert struct {
	ID        string `json:"id"`
	AgendaID  string `json:"agenda_id"`
	Target    string `json:"target"`    // e.g., email address
	Condition string `json:"condition"` // e.g., "always", "room_change"
}
