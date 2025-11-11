package models

type Agenda struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type APIError struct {
	Message string `json:"message"`
}