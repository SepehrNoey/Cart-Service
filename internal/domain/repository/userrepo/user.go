package userrepo

import (
	"context"
	"errors"
	"time"

	"github.com/SepehrNoey/Cart-Service/internal/domain/model"
)

var ErrUserIDDuplicate = errors.New("user id already exists")
var ErrUsernameDuplicate = errors.New("username already exists")
var ErrInvalidCredentials = errors.New("invalid credentials")

type AuthCommand struct {
	UserID    *uint64
	Username  *string
	Password  *string
	Token     *string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

type Repository interface {
	SignUp(ctx context.Context, cmd AuthCommand) error
	Login(ctx context.Context, cmd AuthCommand) (string, error) // return the token and error
	GetUsers(ctx context.Context, cmd AuthCommand) []model.User
}
