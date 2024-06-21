package services

import (
	"context"
	"github.com/iota-agency/iota-erp/internal/domain/position"
)

type PositionService struct {
	repo position.Repository
	app  *Application
}

func NewPositionService(repo position.Repository, app *Application) *PositionService {
	return &PositionService{
		repo: repo,
		app:  app,
	}
}

func (s *PositionService) Count(ctx context.Context) (int64, error) {
	return s.repo.Count(ctx)
}

func (s *PositionService) GetAll(ctx context.Context) ([]*position.Position, error) {
	return s.repo.GetAll(ctx)
}

func (s *PositionService) GetByID(ctx context.Context, id int64) (*position.Position, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *PositionService) GetUploadsPaginated(ctx context.Context, limit, offset int, sortBy []string) ([]*position.Position, error) {
	return s.repo.GetPaginated(ctx, limit, offset, sortBy)
}

func (s *PositionService) Create(ctx context.Context, data *position.Position) error {
	if err := s.repo.Create(ctx, data); err != nil {
		return err
	}
	s.app.EventPublisher.Publish("position.created", data)
	return nil
}

func (s *PositionService) Update(ctx context.Context, data *position.Position) error {
	if err := s.repo.Update(ctx, data); err != nil {
		return err
	}
	s.app.EventPublisher.Publish("position.updated", data)
	return nil
}

func (s *PositionService) Delete(ctx context.Context, id int64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	s.app.EventPublisher.Publish("position.deleted", id)
	return nil
}