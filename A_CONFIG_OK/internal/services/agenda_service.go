package services

import (
	"context"
	"errors"

        "middleware/example/internal/models"
	"middleware/example/internal/repositories"
)

type AgendaService interface {
	List(ctx context.Context) ([]models.Agenda, error)
	Get(ctx context.Context, id string) (*models.Agenda, error)
	Create(ctx context.Context, a models.Agenda) error
	Update(ctx context.Context, a models.Agenda) error
	Delete(ctx context.Context, id string) error
}

type agendaService struct{ repo repositories.AgendaRepository }

func NewAgendaService(r repositories.AgendaRepository) AgendaService { return &agendaService{repo: r} }

func (s *agendaService) List(ctx context.Context) ([]models.Agenda, error) {
	return s.repo.List(ctx)
}

func (s *agendaService) Get(ctx context.Context, id string) (*models.Agenda, error) {
	if id == "" { return nil, errors.New("id required") }
	return s.repo.Get(ctx, id)
}

func (s *agendaService) Create(ctx context.Context, a models.Agenda) error {
	if a.ID == "" || a.Name == "" { return errors.New("id and name are required") }
	return s.repo.Create(ctx, a)
}

func (s *agendaService) Update(ctx context.Context, a models.Agenda) error {
	if a.ID == "" || a.Name == "" { return errors.New("id and name are required") }
	return s.repo.Update(ctx, a)
}

func (s *agendaService) Delete(ctx context.Context, id string) error {
	if id == "" { return errors.New("id required") }
	return s.repo.Delete(ctx, id)
}