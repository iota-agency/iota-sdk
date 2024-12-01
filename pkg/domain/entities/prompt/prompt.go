package prompt

import (
	"time"

	model "github.com/iota-agency/iota-sdk/pkg/interfaces/graph/gqlmodels"
)

type Prompt struct {
	ID          string
	Title       string
	Description string
	Prompt      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (r *Prompt) ToGraph() *model.Prompt {
	return &model.Prompt{
		ID:          r.ID,
		Title:       r.Title,
		Description: r.Description,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}
