package persistence_test

import (
	"context"
	"github.com/gabriel-vasile/mimetype"
	"github.com/iota-uz/iota-sdk/modules/core/domain/entities/upload"
	corepersistence "github.com/iota-uz/iota-sdk/modules/core/infrastructure/persistence"
	persistence2 "github.com/iota-uz/iota-sdk/modules/warehouse/infrastructure/persistence"
	"github.com/jackc/pgx/v5"
	"testing"
	"time"

	"github.com/iota-uz/iota-sdk/modules/warehouse/domain/aggregates/position"
	"github.com/iota-uz/iota-sdk/modules/warehouse/domain/entities/unit"
	"github.com/iota-uz/iota-sdk/pkg/testutils"
)

func TestGormPositionRepository_CRUD(t *testing.T) {
	ctx := testutils.GetTestContext()
	defer func(Tx pgx.Tx, ctx context.Context) {
		err := Tx.Commit(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}(ctx.Tx, ctx.Context)

	unitRepository := persistence2.NewUnitRepository()
	positionRepository := persistence2.NewPositionRepository()
	uploadRepository := corepersistence.NewUploadRepository()

	if err := unitRepository.Create(
		ctx.Context, &unit.Unit{
			ID:         1,
			Title:      "Unit 1",
			ShortTitle: "U1",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}); err != nil {
		t.Fatal(err)
	}

	if err := uploadRepository.Create(
		ctx.Context, &upload.Upload{
			ID:        1,
			Hash:      "hash",
			Path:      "url",
			Size:      1,
			Mimetype:  *mimetype.Lookup("image/png"),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}); err != nil {
		t.Fatal(err)
	}

	if err := positionRepository.Create(
		ctx.Context, &position.Position{
			ID:        1,
			Title:     "Position 1",
			Barcode:   "3141592653589",
			UnitID:    1,
			Images:    []upload.Upload{{ID: 1}},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}); err != nil {
		t.Fatal(err)
	}

	t.Run(
		"GetByID", func(t *testing.T) {
			positionEntity, err := positionRepository.GetByID(ctx.Context, 1)
			if err != nil {
				t.Fatal(err)
			}
			if positionEntity.Title != "Position 1" {
				t.Errorf("expected %s, got %s", "Position 1", positionEntity.Title)
			}
			if positionEntity.Barcode != "3141592653589" {
				t.Errorf("expected %s, got %s", "3141592653589", positionEntity.Barcode)
			}
		},
	)

	t.Run(
		"Update", func(t *testing.T) {
			if err := positionRepository.Update(
				ctx.Context, &position.Position{
					ID:      1,
					Title:   "Updated Position 1",
					Barcode: "3141592653589",
				},
			); err != nil {
				t.Fatal(err)
			}
			positionEntity, err := positionRepository.GetByID(ctx.Context, 1)
			if err != nil {
				t.Fatal(err)
			}
			if positionEntity.Title != "Updated Position 1" {
				t.Errorf("expected %s, got %s", "Updated Position 1", positionEntity.Title)
			}
		},
	)

	t.Run(
		"Delete", func(t *testing.T) {
			if err := positionRepository.Delete(ctx.Context, 1); err != nil {
				t.Fatal(err)
			}
			_, err := positionRepository.GetByID(ctx.Context, 1)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
		},
	)
}
