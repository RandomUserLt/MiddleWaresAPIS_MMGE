package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"middleware/example/internal/models"
	"middleware/example/internal/services"
)

type AgendaController struct{ svc services.AgendaService }

func NewAgendaController(s services.AgendaService) *AgendaController { return &AgendaController{svc: s} }

func (c *AgendaController) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", c.list)
	r.Get("/{id}", c.get)
	r.Post("/", c.create)
	r.Put("/{id}", c.update)
	r.Delete("/{id}", c.delete)
	return r
}

// list
// @Summary      List agendas
// @Tags         agendas
// @Produce      json
// @Success      200  {array}   models.Agenda
// @Failure      500  {object}  models.APIError
// @Router       /agendas [get]
func (c *AgendaController) list(w http.ResponseWriter, r *http.Request) {
	items, err := c.svc.List(r.Context())
	if err != nil { writeErr(w, http.StatusInternalServerError, err.Error()); return }
	writeJSON(w, http.StatusOK, items)
}

// get
// @Summary      Get agenda by id
// @Tags         agendas
// @Param        id   path      string  true  "Agenda ID"
// @Produce      json
// @Success      200  {object}  models.Agenda
// @Failure      404  {object}  models.APIError
// @Router       /agendas/{id} [get]
func (c *AgendaController) get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	a, err := c.svc.Get(r.Context(), id)
	if err != nil { writeErr(w, http.StatusBadRequest, err.Error()); return }
	if a == nil { writeErr(w, http.StatusNotFound, "not found"); return }
	writeJSON(w, http.StatusOK, a)
}

// create
// @Summary      Create agenda
// @Tags         agendas
// @Accept       json
// @Produce      json
// @Param        payload  body  models.Agenda  true  "Agenda"
// @Success      201  {object}  models.Agenda
// @Failure      400  {object}  models.APIError
// @Router       /agendas [post]
func (c *AgendaController) create(w http.ResponseWriter, r *http.Request) {
	var a models.Agenda
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid json"); return
	}
	if err := c.svc.Create(r.Context(), a); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error()); return
	}
	writeJSON(w, http.StatusCreated, a)
}

// update
// @Summary      Update agenda
// @Tags         agendas
// @Accept       json
// @Produce      json
// @Param        id       path  string         true  "Agenda ID"
// @Param        payload  body  models.Agenda  true  "Agenda"
// @Success      200  {object}  models.Agenda
// @Failure      404  {object}  models.APIError
// @Router       /agendas/{id} [put]
func (c *AgendaController) update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var payload models.Agenda
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid json"); return
	}
	payload.ID = id
	if err := c.svc.Update(r.Context(), payload); err != nil {
		if err == sql.ErrNoRows { writeErr(w, http.StatusNotFound, "not found"); return }
		writeErr(w, http.StatusBadRequest, err.Error()); return
	}
	writeJSON(w, http.StatusOK, payload)
}

// delete
// @Summary      Delete agenda
// @Tags         agendas
// @Param        id  path  string  true  "Agenda ID"
// @Success      204  "No Content"
// @Failure      404  {object}  models.APIError
// @Router       /agendas/{id} [delete]
func (c *AgendaController) delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := c.svc.Delete(r.Context(), id); err != nil {
		if err == sql.ErrNoRows { writeErr(w, http.StatusNotFound, "not found"); return }
		writeErr(w, http.StatusBadRequest, err.Error()); return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
func writeErr(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, models.APIError{Message: msg})
}
