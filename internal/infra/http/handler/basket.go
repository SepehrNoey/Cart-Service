package handler

import (
	"errors"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/SepehrNoey/Cart-Service/internal/domain/model"
	"github.com/SepehrNoey/Cart-Service/internal/domain/repository/basketrepo"
	"github.com/SepehrNoey/Cart-Service/internal/infra/http/request"
	"github.com/labstack/echo/v4"
)

const (
	Pending   string = "PENDING"
	Completed string = "COMPLETED"
)

type BasketHandler struct {
	repo basketrepo.Repository
}

func NewBasketHandler(repo basketrepo.Repository) *BasketHandler {
	return &BasketHandler{
		repo: repo,
	}
}

func (bh *BasketHandler) Get(c echo.Context) error {
	baskets := bh.repo.Get(c.Request().Context(), basketrepo.GetCommand{
		ID:        nil,
		CreatedAt: nil,
		UpdatedAt: nil,
		State:     nil,
		Data:      nil,
	})

	if len(baskets) == 0 {
		return echo.ErrNotFound
	}

	return c.JSON(http.StatusOK, baskets)
}

func (bh *BasketHandler) Create(c echo.Context) error {
	var req request.BasketCreate

	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}
	if err := req.CreateValidate(); err != nil {
		return echo.ErrBadRequest // maybe return validation error
	}

	// now, we create a new basket
	id := rand.Uint64() % 1_000_000
	now := time.Now()
	if err := bh.repo.Create(c.Request().Context(), model.Basket{
		ID:        id,
		CreatedAt: now,
		UpdatedAt: now,
		Data:      []byte(req.Data),
		State:     Pending,
	}); err != nil {
		if errors.Is(err, basketrepo.ErrBasketIDDuplicate) {
			return echo.ErrBadRequest
		}
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusCreated, id)
}

func (bh *BasketHandler) Update(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return echo.ErrBadRequest
	}

	var req request.BasketUpdate
	if err = c.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}

	baskets := bh.repo.Get(c.Request().Context(), basketrepo.GetCommand{ID: &id})
	if len(baskets) == 0 {
		return echo.ErrNotFound
	}

	if len(baskets) > 1 {
		return echo.ErrInternalServerError
	}

	// we have found the basket to update
	basket := baskets[0]

	// now, validation
	if err = req.UpdateValidate(); err != nil {
		return echo.ErrBadRequest // maybe return validation error
	}

	if err = bh.repo.Update(c.Request().Context(), model.Basket{
		ID:        basket.ID,
		CreatedAt: basket.CreatedAt,
		UpdatedAt: time.Now(),
		Data:      []byte(req.Data),
		State:     req.State,
	}); err != nil {
		return echo.ErrInternalServerError
	}

	if req.State == Completed {
		// TODO
	}

	return c.JSON(http.StatusOK, basket)
}

func (bh *BasketHandler) GetByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return echo.ErrBadRequest
	}

	baskets := bh.repo.Get(c.Request().Context(), basketrepo.GetCommand{ID: &id})
	if len(baskets) == 0 {
		return echo.ErrNotFound
	}

	if len(baskets) > 1 {
		// must be fixed if happended
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, baskets[0])

}

func (bh *BasketHandler) Delete(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return echo.ErrBadRequest
	}

	baskets := bh.repo.Get(c.Request().Context(), basketrepo.GetCommand{ID: &id})
	if len(baskets) == 0 {
		return echo.ErrNotFound
	}

	if len(baskets) > 1 {
		// must be fixed if happened
		return echo.ErrInternalServerError
	}

	if err = bh.repo.Delete(c.Request().Context(), basketrepo.GetCommand{ID: &id}); err != nil {
		return echo.ErrInternalServerError
	}

	// get all baskets again
	baskets = bh.repo.Get(c.Request().Context(), basketrepo.GetCommand{})
	return c.JSON(http.StatusOK, baskets)

}

func (bh *BasketHandler) Register(g *echo.Group) {
	g.GET("", bh.Get)
	g.POST("", bh.Create)
	g.PATCH(":id", bh.Update)
	g.GET(":id", bh.GetByID)
	g.DELETE(":id", bh.Delete)
}
