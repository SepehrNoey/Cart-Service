package basketrepo

import (
	"context"
	"errors"
	"time"

	"github.com/SepehrNoey/Cart-Service/internal/domain/model"
)

var ErrBasketIDDuplicate = errors.New("given basket id already exists")

type GetCommand struct {
	ID        *uint64
	CreatedAt *time.Time
	UpdatedAt *time.Time
	Data      *[]byte
	State     *string
}

type Repository interface {
	Create(ctx context.Context, basket model.Basket) error
	Get(ctx context.Context, cmd GetCommand) []model.Basket
	Update(ctx context.Context, cmd GetCommand, basket model.Basket) error
	Delete(ctx context.Context, cmd GetCommand) error
}
