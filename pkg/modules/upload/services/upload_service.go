package services

import (
	"context"

	"github.com/iota-agency/iota-sdk/pkg/composables"
	"github.com/iota-agency/iota-sdk/pkg/event"
	"github.com/iota-agency/iota-sdk/pkg/modules/upload/domain/entities/upload"
	"github.com/iota-agency/iota-sdk/pkg/modules/upload/permissions"
)

type UploadService struct {
	repo      upload.Repository
	publisher event.Publisher
}

func NewUploadService(
	repo upload.Repository,
	publisher event.Publisher,
) *UploadService {
	return &UploadService{
		repo:      repo,
		publisher: publisher,
	}
}

func (s *UploadService) GetByID(ctx context.Context, id string) (*upload.Upload, error) {
	if err := composables.CanUser(ctx, permissions.UploadRead); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *UploadService) GetAll(ctx context.Context) ([]*upload.Upload, error) {
	if err := composables.CanUser(ctx, permissions.UploadRead); err != nil {
		return nil, err
	}
	return s.repo.GetAll(ctx)
}

func (s *UploadService) Create(ctx context.Context, data *upload.CreateDTO) error {
	if err := composables.CanUser(ctx, permissions.UploadCreate); err != nil {
		return err
	}
	entity, err := data.ToEntity()
	if err != nil {
		return err
	}
	if err := s.repo.Create(ctx, entity); err != nil {
		return err
	}
	createdEvent, err := upload.NewCreatedEvent(ctx, *data, *entity)
	if err != nil {
		return err
	}
	s.publisher.Publish(createdEvent)
	return nil
}

func (s *UploadService) Update(ctx context.Context, id string, data *upload.UpdateDTO) error {
	if err := composables.CanUser(ctx, permissions.UploadUpdate); err != nil {
		return err
	}
	entity, err := data.ToEntity(id)
	if err != nil {
		return err
	}
	if err := s.repo.Update(ctx, entity); err != nil {
		return err
	}
	updatedEvent, err := upload.NewUpdatedEvent(ctx, *data, *entity)
	if err != nil {
		return err
	}
	s.publisher.Publish(updatedEvent)
	return nil
}

func (s *UploadService) Delete(ctx context.Context, id string) (*upload.Upload, error) {
	if err := composables.CanUser(ctx, permissions.UploadDelete); err != nil {
		return nil, err
	}
	entity, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return nil, err
	}
	deletedEvent, err := upload.NewDeletedEvent(ctx, *entity)
	if err != nil {
		return nil, err
	}
	s.publisher.Publish(deletedEvent)
	return entity, nil
}
