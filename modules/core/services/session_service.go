package services

import (
	"context"
	"github.com/iota-agency/iota-sdk/pkg/domain/entities/session"
	"github.com/iota-agency/iota-sdk/pkg/event"
)

type SessionService struct {
	repo      session.Repository
	publisher event.Publisher
}

func NewSessionService(repo session.Repository, publisher event.Publisher) *SessionService {
	return &SessionService{
		repo:      repo,
		publisher: publisher,
	}
}

func (s *SessionService) GetCount(ctx context.Context) (int64, error) {
	return s.repo.Count(ctx)
}

func (s *SessionService) GetAll(ctx context.Context) ([]*session.Session, error) {
	return s.repo.GetAll(ctx)
}

func (s *SessionService) GetByToken(ctx context.Context, id string) (*session.Session, error) {
	return s.repo.GetByToken(ctx, id)
}

func (s *SessionService) GetPaginated(
	ctx context.Context,
	limit, offset int,
	sortBy []string,
) ([]*session.Session, error) {
	return s.repo.GetPaginated(ctx, limit, offset, sortBy)
}

func (s *SessionService) Create(ctx context.Context, data *session.CreateDTO) error {
	entity := data.ToEntity()
	if err := s.repo.Create(ctx, entity); err != nil {
		return err
	}
	createdEvent, err := session.NewCreatedEvent(*data, *entity)
	if err != nil {
		return err
	}
	s.publisher.Publish(createdEvent)
	return nil
}

func (s *SessionService) Update(ctx context.Context, data *session.Session) error {
	if err := s.repo.Update(ctx, data); err != nil {
		return err
	}
	s.publisher.Publish("session.updated", data)
	return nil
}

func (s *SessionService) Delete(ctx context.Context, id int64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	s.publisher.Publish("session.deleted", id)
	return nil
}
