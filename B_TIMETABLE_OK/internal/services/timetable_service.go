package services

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	ics "github.com/arran4/golang-ical"
	"middleware/example/internal/models"
)

type TimetableService interface {
	FetchEvents(agendaIDs []string, from, to *time.Time) ([]models.Event, error)
}

type timetableService struct {
	client     *http.Client
	IcalBase   string // base URL without resources=...
	WeeksParam string // nbWeeks to request (string form)
}

func NewTimetableService(client *http.Client) TimetableService {
	// Defaults from the TP spec
	return &timetableService{
		client:     client,
		IcalBase:   "https://edt.uca.fr/jsp/custom/modules/plannings/anonymous_cal.jsp?projectId=3&calType=ical&displayConfigId=128",
		WeeksParam: "8",
	}
}

func (s *timetableService) buildURL(agendaIDs []string) (string, error) {
	if len(agendaIDs) == 0 {
		return "", errors.New("agendaIds required")
	}
	joined := strings.Join(agendaIDs, ",")
	return fmt.Sprintf("%s&nbWeeks=%s&resources=%s", s.IcalBase, s.WeeksParam, joined), nil
}

func parseICalTime(v string) (time.Time, error) {
	// Try typical iCal formats (Zulu, with offset, or naive)
	layouts := []string{
		"20060102T150405Z",
		"20060102T150405-0700",
		"20060102T150405",
		"20060102",
	}
	var last error
	for _, layout := range layouts {
		if t, err := time.Parse(layout, v); err == nil {
			return t, nil
		} else {
			last = err
		}
	}
	return time.Time{}, last
}

func toRFC3339Ptr(t *time.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}
	return t.UTC().Format(time.RFC3339)
}

func (s *timetableService) FetchEvents(agendaIDs []string, from, to *time.Time) ([]models.Event, error) {
	url, err := s.buildURL(agendaIDs)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return nil, fmt.Errorf("ical fetch %d: %s", resp.StatusCode, string(b))
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	cal, err := ics.ParseCalendar(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	var out []models.Event
	for _, comp := range cal.Components {
		ve, ok := comp.(*ics.VEvent)
		if !ok {
			continue
		}
		uid := ve.GetProperty("UID")
		if uid == nil || uid.Value == "" {
			continue
		}

		// DTSTART / DTEND
		var start, end time.Time
		if p := ve.GetProperty("DTSTART"); p != nil {
			if ts, err := parseICalTime(p.Value); err == nil {
				start = ts
			}
		}
		if p := ve.GetProperty("DTEND"); p != nil {
			if te, err := parseICalTime(p.Value); err == nil {
				end = te
			}
		}

		// Filtering window [from, to]
		if from != nil && !from.IsZero() && end.Before(from.UTC()) {
			continue
		}
		if to != nil && !to.IsZero() && start.After(to.UTC()) {
			continue
		}

		title := ""
		if p := ve.GetProperty("SUMMARY"); p != nil {
			title = p.Value
		}
		location := ""
		if p := ve.GetProperty("LOCATION"); p != nil {
			location = p.Value
		}
		description := ""
		if p := ve.GetProperty("DESCRIPTION"); p != nil {
			description = strings.ReplaceAll(p.Value, "\\n", "\n")
		}
		lastMod := ""
		if p := ve.GetProperty("LAST-MODIFIED"); p != nil {
			if tm, err := parseICalTime(p.Value); err == nil {
				lastMod = tm.UTC().Format(time.RFC3339)
			}
		}

		out = append(out, models.Event{
			ID:          uid.Value,
			AgendaIDs:   agendaIDs,
			Title:       title,
			Description: description,
			Start:       start.UTC().Format(time.RFC3339),
			End:         end.UTC().Format(time.RFC3339),
			Location:    location,
			LastUpdate:  lastMod,
		})
	}
	if out == nil {
		out = []models.Event{}
	}
	return out, nil
}
