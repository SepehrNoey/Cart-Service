package basketrepo

import (
	"context"
	"errors"
	"time"

	"github.com/SepehrNoey/Cart-Service/internal/domain/model"
	"github.com/SepehrNoey/Cart-Service/internal/infra/repository/jsonb"
)

var ErrBasketIDDuplicate = errors.New("given basket id already exists")

type GetCommand struct {
	ID        *uint64
	CreatedAt *time.Time
	UpdatedAt *time.Time
	Data      *string
	State     *string
}

type Repository interface {
	Create(ctx context.Context, basket model.Basket) error
	Get(ctx context.Context, cmd GetCommand) ([]model.Basket, []jsonb.JSONB)
	Update(ctx context.Context, basket model.Basket) error
	Delete(ctx context.Context, cmd GetCommand) error
}
