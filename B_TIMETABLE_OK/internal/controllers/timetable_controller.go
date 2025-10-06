package controllers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
        
	"github.com/go-chi/chi/v5"
	"middleware/example/internal/models"
	"middleware/example/internal/services"
)

type TimetableController struct{ svc services.TimetableService }

func NewTimetableController(s services.TimetableService) *TimetableController { return &TimetableController{svc: s} }






// @Summary      List events from UCA iCal
// @Tags         timetable
// @Produce      json
// @Param        agendaIds  query  string  true  "Comma-separated agenda IDs (e.g., 13295,13345)"
// @Param        from       query  string  false "ISO date (YYYY-MM-DD) start filter"
// @Param        to         query  string  false "ISO date (YYYY-MM-DD) end filter"
// @Success      200  {array}   models.Event
// @Failure      400  {object}  models.APIError
// @Failure      502  {object}  models.APIError
// @Router       /events [get]
func (c *TimetableController) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	rawIDs := strings.TrimSpace(q.Get("agendaIds"))
	if rawIDs == "" {
		writeJSON(w, http.StatusBadRequest, models.APIError{Message: "agendaIds is required"})
		return
	}
	var agendaIDs []string
	for _, s := range strings.Split(rawIDs, ",") {
		s = strings.TrimSpace(s)
		if s != "null" && s != "" {
			agendaIDs = append(agendaIDs, s)
		}
	}
	if len(agendaIDs) == 0 {
		writeJSON(w, http.StatusBadRequest, models.APIError{Message: "agendaIds is required"})
		return
	}

	var fromPtr, toPtr *time.Time
	if v := strings.TrimSpace(q.Get("from")); v != "" {
		if t, err := time.Parse("2006-01-02", v); err == nil {
			ft := t.UTC()
			fromPtr = &ft
		} else {
			writeJSON(w, http.StatusBadRequest, models.APIError{Message: "invalid from (YYYY-MM-DD)"})
			return
		}
	}
	if v := strings.TrimSpace(q.Get("to")); v != "" {
		if t, err := time.Parse("2006-01-02", v); err == nil {
			tt := t.Add(23*time.Hour + 59*time.Minute + 59*time.Second).UTC()
			toPtr = &tt
		} else {
			writeJSON(w, http.StatusBadRequest, models.APIError{Message: "invalid to (YYYY-MM-DD)"})
			return
		}
	}

	events, err := c.svc.FetchEvents(agendaIDs, fromPtr, toPtr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		_ = json.NewEncoder(w).Encode(models.APIError{Message: err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if events == nil {
		events = []models.Event{}
	}
	_ = json.NewEncoder(w).Encode(events)
}



func (c *TimetableController) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, `{"message":"invalid id"}`, http.StatusBadRequest)
		return
	}

	agendaIdsParam := strings.TrimSpace(r.URL.Query().Get("agendaIds"))
	if agendaIdsParam == "" {
		http.Error(w, `{"message":"agendaIds required"}`, http.StatusBadRequest)
		return
	}
	agendaIDs := strings.Split(agendaIdsParam, ",")

	fromPtr, toPtr, err := parseFromTo(r)
	if err != nil {
		http.Error(w, `{"message":"invalid date format (use YYYY-MM-DD)"}`, http.StatusBadRequest)
		return
	}

	events, err := c.svc.FetchEvents(agendaIDs, fromPtr, toPtr)
	if err != nil {
		http.Error(w, `{"message":"ical fetch failed"}`, http.StatusBadGateway)
		return
	}

	for _, ev := range events {
		if ev.ID == id {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(ev)
			return
		}
	}

	http.Error(w, `{"message":"not found"}`, http.StatusNotFound)
}

// parseFromTo lit les query params from/to (format YYYY-MM-DD) et renvoie des *time.Time
func parseFromTo(r *http.Request) (fromPtr, toPtr *time.Time, err error) {
	q := r.URL.Query()

	if v := strings.TrimSpace(q.Get("from")); v != "" {
		t, e := time.Parse("2006-01-02", v)
		if e != nil {
			return nil, nil, e
		}
		ft := t.UTC()
		fromPtr = &ft
	}

	if v := strings.TrimSpace(q.Get("to")); v != "" {
		t, e := time.Parse("2006-01-02", v)
		if e != nil {
			return nil, nil, e
		}
		// inclure toute la journée "to"
		tt := t.Add(23*time.Hour + 59*time.Minute + 59*time.Second).UTC()
		toPtr = &tt
	}

	return fromPtr, toPtr, nil
}

