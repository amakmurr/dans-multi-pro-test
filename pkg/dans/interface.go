package dans

import (
	"context"

	"github.com/google/uuid"
)

type ClientInterface interface {
	GetJobList(ctx context.Context, params *GetJobListParams) ([]*Job, error)
	GetJobDetail(ctx context.Context, id uuid.UUID) (*Job, error)
}
