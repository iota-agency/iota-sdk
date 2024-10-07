package authlog

import (
	"time"

	model "github.com/iota-agency/iota-erp/internal/interfaces/graph/gqlmodels"
)

type AuthenticationLog struct {
	ID        int64
	UserID    int64
	IP        string
	UserAgent string
	CreatedAt time.Time
}

func (r *AuthenticationLog) ToGraph() *model.AuthenticationLog {
	return &model.AuthenticationLog{
		ID:        r.ID,
		UserID:    r.UserID,
		IP:        r.IP,
		UserAgent: r.UserAgent,
		CreatedAt: r.CreatedAt,
	}
}
