package basketrepo

import (
	"context"
	"errors"
	"time"

	"github.com/SepehrNoey/Cart-Service/internal/domain/model"
	"github.com/SepehrNoey/Cart-Service/internal/infra/repository/jsonb"
)

var ErrBasketIDDuplicate = errors.New("given basket id already exists")
var ErrCompletedBasketCantChange = errors.New("completed basket can't change")
var ErrBasketDataInvalidLength = errors.New("invalid data length, at most 2048 bytes")

type GetCommand struct {
	ID        *uint64
	CreatedAt *time.Time
	UpdatedAt *time.Time
	Data      *string
	State     *string
	UserID    *uint64
}

type Repository interface {
	Create(ctx context.Context, basket model.Basket) error
	Get(ctx context.Context, cmd GetCommand) ([]model.Basket, []jsonb.JSONB)
	Update(ctx context.Context, basket model.Basket) error
	Delete(ctx context.Context, cmd GetCommand) error
}
