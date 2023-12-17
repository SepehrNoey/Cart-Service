package usersql

import (
	"context"
	"errors"
	"time"

	"github.com/SepehrNoey/Cart-Service/internal/domain/model"
	"github.com/SepehrNoey/Cart-Service/internal/domain/repository/userrepo"
	"github.com/SepehrNoey/Cart-Service/internal/infra/http/auth"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type UserDTO struct {
	model.User
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	token     string
}

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) SignUp(ctx context.Context, cmd userrepo.AuthCommand) error {
	if cmd.Username == nil || cmd.Password == nil {
		return userrepo.ErrInvalidCredentials
	}

	users := r.GetUsers(ctx, userrepo.AuthCommand{Username: cmd.Username})
	if len(users) > 0 {
		return userrepo.ErrUsernameDuplicate
	}

	now := time.Now()
	dto := UserDTO{User: model.User{ID: *cmd.UserID, Username: *cmd.Username, Password: *cmd.Password},
		CreatedAt: now, UpdatedAt: now}
	if result := r.db.Model(&UserDTO{}).WithContext(ctx).Create(&dto); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return userrepo.ErrUserIDDuplicate
		}
	}
	return nil

}

// the AuthCommand must be set carefully before calling this function
func (r *Repository) Login(ctx context.Context, cmd userrepo.AuthCommand) (string, error) {
	if cmd.Username == nil || cmd.Password == nil {
		return "", userrepo.ErrInvalidCredentials
	}

	users := r.GetUsers(ctx, cmd)
	if len(users) == 0 {
		return "", userrepo.ErrInvalidCredentials
	}
	if len(users) > 1 {
		return "", echo.ErrInternalServerError
	}

	user := users[0]
	token, err := auth.CreateToken(user.ID, user.Username)
	if err != nil {
		return "", echo.ErrInternalServerError
	}

	return token, nil
}

func (r *Repository) GetUsers(ctx context.Context, cmd userrepo.AuthCommand) []model.User {
	var userDTOs []UserDTO

	var dto UserDTO
	var conditions []string

	if cmd.UserID != nil {
		dto.ID = *cmd.UserID
		conditions = append(conditions, "ID")
	}
	if cmd.Username != nil {
		dto.Username = *cmd.Username
		conditions = append(conditions, "Username")
	}
	if cmd.Password != nil {
		dto.Password = *cmd.Password
		conditions = append(conditions, "Password")
	}
	if cmd.Token != nil {
		dto.token = *cmd.Token
		conditions = append(conditions, "Token")
	}
	if cmd.CreatedAt != nil {
		dto.CreatedAt = *cmd.CreatedAt
		conditions = append(conditions, "CreatedAt")
	}
	if cmd.UpdatedAt != nil {
		dto.UpdatedAt = *cmd.UpdatedAt
		conditions = append(conditions, "UpdatedAt")
	}

	if len(conditions) == 0 {
		if err := r.db.WithContext(ctx).Find(&userDTOs); err.Error != nil {
			return nil // maybe should change here to return and error
		}
	} else {
		if err := r.db.WithContext(ctx).Where(&dto, conditions).Find(&userDTOs); err.Error != nil {
			return nil
		}
	}

	users := make([]model.User, len(userDTOs))
	for i, dto := range userDTOs {
		users[i] = dto.User
	}

	return users
}
