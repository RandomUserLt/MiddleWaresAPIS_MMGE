package models

// Reçu depuis le consumer Timetable (message sur NATS)
type Change struct {
	Field string `json:"field"` // location,start,end,title,description
	Old   string `json:"old"`
	New   string `json:"new"`
	Kind  string `json:"kind"` // room_change,time_change, etc.
}

type AlertEvent struct {
	Type       string   `json:"type"`       // event_changed | event_new
	EventID    string   `json:"eventId"`
	AgendaIDs  []string `json:"agendaIds"`
	Title      string   `json:"title"`
	Start      string   `json:"start"`
	End        string   `json:"end"`
	Location   string   `json:"location"`
	Changes    []Change `json:"changes"`
	EmailText  string   `json:"emailText"`  // déjà prêt, utile en fallback
}

// Récupérés depuis l’API Config : qui doit recevoir quoi
type AlertSubscription struct {
	ID        string `json:"id"`
	AgendaID  string `json:"agenda_id"` // "all" possible si tu l’as prévu
	Target    string `json:"target"`    // adresse email
	Condition string `json:"condition"` // always | room_change | time_change | ...
}

// Requête d’envoi à l’API mail GCC
//type OutgoingMail struct {
//	To      string `json:"to"`      // destinataire
//	Subject string `json:"subject"` // objet
//	Text    string `json:"text"`    // corps en texte simple (ou HTML, cf. doc API)
//}
//new
type OutgoingMail struct {
    Recipient string `json:"recipient"`
    Subject   string `json:"subject"`
    Content   string `json:"content"`
}









