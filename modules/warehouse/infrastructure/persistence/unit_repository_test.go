package persistence_test

import (
	"context"
	"github.com/iota-uz/iota-sdk/modules/warehouse/domain/entities/unit"
	"github.com/iota-uz/iota-sdk/modules/warehouse/infrastructure/persistence"
	"github.com/iota-uz/iota-sdk/pkg/testutils"
	"github.com/jackc/pgx/v5"
	"log"
	"testing"
	"time"
)

func TestGormUnitRepository_CRUD(t *testing.T) {
	ctx := testutils.GetTestContext()
	defer func(Tx pgx.Tx, ctx context.Context) {
		err := Tx.Commit(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(ctx.Tx, ctx.Context)
	unitRepository := persistence.NewUnitRepository()

	if err := unitRepository.Create(
		ctx.Context, &unit.Unit{
			ID:         1,
			Title:      "test",
			ShortTitle: "t",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	); err != nil {
		t.Fatal(err)
	}

	t.Run(
		"Count", func(t *testing.T) {
			count, err := unitRepository.Count(ctx.Context)
			if err != nil {
				t.Fatal(err)
			}
			if count != 1 {
				t.Errorf("expected 1, got %d", count)
			}
		},
	)

	t.Run(
		"GetPaginated", func(t *testing.T) {
			accounts, err := unitRepository.GetPaginated(ctx.Context, &unit.FindParams{Limit: 1})
			if err != nil {
				t.Fatal(err)
			}
			if len(accounts) != 1 {
				t.Errorf("expected 1, got %d", len(accounts))
			}
		},
	)

	t.Run(
		"GetAll", func(t *testing.T) {
			units, err := unitRepository.GetAll(ctx.Context)
			if err != nil {
				t.Fatal(err)
			}
			if len(units) != 1 {
				t.Errorf("expected 1, got %d", len(units))
			}
			if units[0].Title != "test" {
				t.Errorf("expected test, got %s", units[0].Title)
			}
		},
	)

	t.Run(
		"GetByID", func(t *testing.T) {
			unitEntity, err := unitRepository.GetByID(ctx.Context, 1)
			if err != nil {
				t.Fatal(err)
			}
			if unitEntity.Title != "test" {
				t.Errorf("expected test, got %s", unitEntity.Title)
			}
		},
	)

	t.Run(
		"Update", func(t *testing.T) {
			if err := unitRepository.Update(
				ctx.Context, &unit.Unit{
					ID:         1,
					Title:      "test2",
					ShortTitle: "t2",
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				},
			); err != nil {
				t.Fatal(err)
			}
			unitEntity, err := unitRepository.GetByID(ctx.Context, 1)
			if err != nil {
				t.Fatal(err)
			}
			if unitEntity.Title != "test2" {
				t.Errorf("expected test2, got %s", unitEntity.Title)
			}
		},
	)
}
