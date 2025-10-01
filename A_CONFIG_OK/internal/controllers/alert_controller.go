package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"middleware/example/internal/models"
	"middleware/example/internal/services"
)

type AlertController struct{ svc services.AlertService }

func NewAlertController(s services.AlertService) *AlertController { return &AlertController{svc: s} }

func (c *AlertController) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", c.list)
	r.Get("/{id}", c.get)
	r.Post("/", c.create)
	r.Put("/{id}", c.update)
	r.Delete("/{id}", c.delete)
	return r
}

// list
// @Summary      List alerts (optional filter by agenda_id)
// @Tags         alerts
// @Produce      json
// @Param        agenda_id  query  string  false  "Filter by agenda id"
// @Success      200  {array}   models.Alert
// @Failure      400  {object}  models.APIError
// @Router       /alerts [get]
func (c *AlertController) list(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("agenda_id")
	var filter *string
	if q != "" {
		filter = &q
	}
	items, err := c.svc.List(r.Context(), filter)
	if err != nil { writeErr(w, http.StatusBadRequest, err.Error()); return }
	writeJSON(w, http.StatusOK, items)
}

// get
// @Summary      Get alert by id
// @Tags         alerts
// @Param        id   path  string  true  "Alert ID"
// @Produce      json
// @Success      200  {object}  models.Alert
// @Failure      404  {object}  models.APIError
// @Router       /alerts/{id} [get]
func (c *AlertController) get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	a, err := c.svc.Get(r.Context(), id)
	if err != nil { writeErr(w, http.StatusBadRequest, err.Error()); return }
	if a == nil { writeErr(w, http.StatusNotFound, "not found"); return }
	writeJSON(w, http.StatusOK, a)
}

// create
// @Summary      Create alert
// @Tags         alerts
// @Accept       json
// @Produce      json
// @Param        payload  body  models.Alert  true  "Alert"
// @Success      201  {object}  models.Alert
// @Failure      400  {object}  models.APIError
// @Router       /alerts [post]
func (c *AlertController) create(w http.ResponseWriter, r *http.Request) {
	var a models.Alert
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid json"); return
	}
	if err := c.svc.Create(r.Context(), a); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error()); return
	}
	writeJSON(w, http.StatusCreated, a)
}

// update
// @Summary      Update alert
// @Tags         alerts
// @Accept       json
// @Produce      json
// @Param        id       path  string        true  "Alert ID"
// @Param        payload  body  models.Alert  true  "Alert"
// @Success      200  {object}  models.Alert
// @Failure      404  {object}  models.APIError
// @Router       /alerts/{id} [put]
func (c *AlertController) update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var payload models.Alert
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
// @Summary      Delete alert
// @Tags         alerts
// @Param        id  path  string  true  "Alert ID"
// @Success      204  "No Content"
// @Failure      404  {object}  models.APIError
// @Router       /alerts/{id} [delete]
func (c *AlertController) delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := c.svc.Delete(r.Context(), id); err != nil {
		if err == sql.ErrNoRows { writeErr(w, http.StatusNotFound, "not found"); return }
		writeErr(w, http.StatusBadRequest, err.Error()); return
	}
	w.WriteHeader(http.StatusNoContent)
}

// NOTE: writeJSON and writeErr are already available in agenda_controller.go (same package).
