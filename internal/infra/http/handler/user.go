package handler

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/SepehrNoey/Cart-Service/internal/domain/repository/userrepo"
	"github.com/SepehrNoey/Cart-Service/internal/infra/http/request"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	repo userrepo.Repository
}

func NewUserHandler(repo userrepo.Repository) *UserHandler {
	return &UserHandler{
		repo: repo,
	}
}

func (uh *UserHandler) SignUp(c echo.Context) error {
	var req request.AuthUser
	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}

	users := uh.repo.GetUsers(c.Request().Context(), userrepo.AuthCommand{Username: &req.Username})
	if len(users) != 0 {
		return userrepo.ErrUsernameDuplicate
	}

	userID := rand.Uint64() % 1_000_000
	now := time.Now()
	if err := uh.repo.SignUp(c.Request().Context(), userrepo.AuthCommand{
		UserID:    &userID,
		Username:  &req.Username,
		Password:  &req.Password,
		CreatedAt: &now,
		UpdatedAt: &now,
	}); err != nil {
		return err
	}

	return c.JSONPretty(http.StatusOK, req.Username, "  ")

}

func (uh *UserHandler) Login(c echo.Context) error {
	var req request.AuthUser
	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}

	users := uh.repo.GetUsers(c.Request().Context(), userrepo.AuthCommand{Username: &req.Username, Password: &req.Password})
	if len(users) == 0 {
		return echo.ErrNotFound
	}
	if len(users) > 1 {
		return echo.ErrInternalServerError
	}
	token, err := uh.repo.Login(c.Request().Context(), userrepo.AuthCommand{Username: &req.Username, Password: &req.Password})
	if err != nil {
		return err
	}

	return c.JSONPretty(http.StatusOK, token, "  ")
}

func (uh *UserHandler) Register(g *echo.Group) {
	g.POST("signup", uh.SignUp)
	g.POST("login", uh.Login)

}
