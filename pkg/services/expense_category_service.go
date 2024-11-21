package services

import (
	"context"
	"github.com/iota-agency/iota-sdk/pkg/event"

	category "github.com/iota-agency/iota-sdk/pkg/domain/aggregates/expense_category"
)

type ExpenseCategoryService struct {
	repo      category.Repository
	publisher event.Publisher
}

func NewExpenseCategoryService(repo category.Repository, publisher event.Publisher) *ExpenseCategoryService {
	return &ExpenseCategoryService{
		repo:      repo,
		publisher: publisher,
	}
}

func (s *ExpenseCategoryService) GetByID(ctx context.Context, id uint) (*category.ExpenseCategory, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ExpenseCategoryService) Count(ctx context.Context) (uint, error) {
	return s.repo.Count(ctx)
}

func (s *ExpenseCategoryService) GetAll(ctx context.Context) ([]*category.ExpenseCategory, error) {
	return s.repo.GetAll(ctx)
}

func (s *ExpenseCategoryService) GetPaginated(
	ctx context.Context,
	limit, offset int,
	sortBy []string,
) ([]*category.ExpenseCategory, error) {
	return s.repo.GetPaginated(ctx, limit, offset, sortBy)
}

func (s *ExpenseCategoryService) Create(ctx context.Context, data *category.CreateDTO) error {
	createdEvent, err := category.NewCreatedEvent(ctx, *data)
	if err != nil {
		return err
	}
	entity, err := data.ToEntity()
	if err != nil {
		return err
	}
	if err := s.repo.Create(ctx, entity); err != nil {
		return err
	}
	createdEvent.Result = *entity
	s.publisher.Publish(createdEvent)
	return nil
}

func (s *ExpenseCategoryService) Update(ctx context.Context, id uint, data *category.UpdateDTO) error {
	updatedEvent, err := category.NewUpdatedEvent(ctx, *data)
	if err != nil {
		return err
	}
	entity, err := data.ToEntity(id)
	if err != nil {
		return err
	}
	if err := s.repo.Update(ctx, entity); err != nil {
		return err
	}
	updatedEvent.Result = *entity
	s.publisher.Publish(updatedEvent)
	return nil
}

func (s *ExpenseCategoryService) Delete(ctx context.Context, id uint) (*category.ExpenseCategory, error) {
	deletedEvent, err := category.NewDeletedEvent(ctx)
	if err != nil {
		return nil, err
	}
	entity, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return nil, err
	}
	deletedEvent.Result = *entity
	s.publisher.Publish(deletedEvent)
	return entity, nil
}