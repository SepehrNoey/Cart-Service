package handler

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/SepehrNoey/Cart-Service/internal/domain/model"
	"github.com/SepehrNoey/Cart-Service/internal/domain/repository/basketrepo"
	"github.com/SepehrNoey/Cart-Service/internal/infra/http/auth"
	"github.com/SepehrNoey/Cart-Service/internal/infra/http/request"
	"github.com/golang-jwt/jwt/v5"
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
	var req request.BasketGet

	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}
	claims, err := auth.ValidateToken(req.Token)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) || errors.Is(err, jwt.ErrTokenInvalidClaims) {
			return echo.ErrUnauthorized
		}
		return echo.ErrInternalServerError
	}

	userID, err := strconv.ParseUint(fmt.Sprintf("%v", claims["user_id"]), 10, 64)
	if err != nil {
		return echo.ErrInternalServerError
	}

	baskets, jsonbs := bh.repo.Get(c.Request().Context(), basketrepo.GetCommand{
		ID:        nil,
		CreatedAt: nil,
		UpdatedAt: nil,
		State:     nil,
		Data:      nil,
		UserID:    &userID,
	})

	if len(baskets) == 0 {
		return echo.ErrNotFound
	}

	resultMap := make(map[string]interface{})
	resultMap["baskets"] = baskets
	resultMap["jsonbs"] = jsonbs
	return c.JSONPretty(http.StatusOK, resultMap, "  ")
}

func (bh *BasketHandler) Create(c echo.Context) error {
	var req request.BasketCreate

	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}

	claims, err := auth.ValidateToken(req.Token)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) || errors.Is(err, jwt.ErrTokenInvalidClaims) {
			return echo.ErrUnauthorized
		}
		return echo.ErrInternalServerError
	}
	if err := req.CreateValidate(); err != nil {
		return err
	}

	userID, err := strconv.ParseUint(fmt.Sprintf("%v", claims["user_id"]), 10, 64)
	if err != nil {
		return echo.ErrInternalServerError
	}

	// now, we create a new basket
	basketID := rand.Uint64() % 1_000_000
	now := time.Now()
	basket := model.Basket{
		ID:        basketID,
		CreatedAt: now,
		UpdatedAt: now,
		Data:      req.Data,
		State:     Pending,
		UserID:    userID,
	}
	if err := bh.repo.Create(c.Request().Context(), basket); err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSONPretty(http.StatusCreated, basket, "  ")
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

	claims, err := auth.ValidateToken(req.Token)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) || errors.Is(err, jwt.ErrTokenInvalidClaims) {
			return echo.ErrUnauthorized
		}
		return echo.ErrInternalServerError
	}

	userID, err := strconv.ParseUint(fmt.Sprintf("%v", claims["user_id"]), 10, 64)
	if err != nil {
		return echo.ErrInternalServerError
	}

	baskets, _ := bh.repo.Get(c.Request().Context(), basketrepo.GetCommand{
		ID:        &id,
		CreatedAt: nil,
		UpdatedAt: nil,
		Data:      nil,
		State:     nil,
		UserID:    &userID,
	})
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
		return err
	}

	var toBeUpdatedData string
	if req.Data == "" {
		toBeUpdatedData = basket.Data
	} else {
		toBeUpdatedData = req.Data
	}

	var toBeUpdatedState string
	if req.State == "" {
		toBeUpdatedState = basket.State
	} else if req.State == Completed && basket.State == Pending {
		toBeUpdatedState = req.State
	} else {
		if req.State == Pending && basket.State == Completed {
			return basketrepo.ErrCompletedBasketCantChange
		}
		toBeUpdatedState = basket.State
	}

	if err = bh.repo.Update(c.Request().Context(), model.Basket{
		ID:        basket.ID,
		CreatedAt: basket.CreatedAt,
		UpdatedAt: time.Now(),
		Data:      toBeUpdatedData,
		State:     toBeUpdatedState,
		UserID:    userID,
	}); err != nil {
		return echo.ErrInternalServerError
	}

	baskets, jsonbs := bh.repo.Get(c.Request().Context(), basketrepo.GetCommand{ID: &basket.ID, UserID: &userID})
	basket = baskets[0]

	resultMap := make(map[string]interface{})
	resultMap["updated_basket"] = basket
	resultMap["jsonb"] = jsonbs[0]
	return c.JSONPretty(http.StatusOK, resultMap, "  ")
}

func (bh *BasketHandler) GetByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return echo.ErrBadRequest
	}

	var req request.BasketGet

	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}
	claims, err := auth.ValidateToken(req.Token)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) || errors.Is(err, jwt.ErrTokenInvalidClaims) {
			return echo.ErrUnauthorized
		}
		return echo.ErrInternalServerError
	}

	userID, err := strconv.ParseUint(fmt.Sprintf("%v", claims["user_id"]), 10, 64) // convert to string first
	if err != nil {
		return echo.ErrInternalServerError
	}

	baskets, jsonbs := bh.repo.Get(c.Request().Context(), basketrepo.GetCommand{ID: &id, UserID: &userID})
	if len(baskets) == 0 {
		return echo.ErrNotFound
	}

	if len(baskets) > 1 {
		// must be fixed if happended
		return echo.ErrInternalServerError
	}

	resultMap := make(map[string]interface{})
	resultMap["basket_by_id"] = baskets[0]
	resultMap["jsonb"] = jsonbs[0]

	return c.JSONPretty(http.StatusOK, resultMap, "  ")

}

func (bh *BasketHandler) Delete(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return echo.ErrBadRequest
	}

	var req request.BasketGet

	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}
	claims, err := auth.ValidateToken(req.Token)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) || errors.Is(err, jwt.ErrTokenInvalidClaims) {
			return echo.ErrUnauthorized
		}
		return echo.ErrInternalServerError
	}

	userID, err := strconv.ParseUint(fmt.Sprintf("%v", claims["user_id"]), 10, 64)
	if err != nil {
		return echo.ErrInternalServerError
	}

	baskets, _ := bh.repo.Get(c.Request().Context(), basketrepo.GetCommand{ID: &id, UserID: &userID})
	if len(baskets) == 0 {
		return echo.ErrNotFound
	}

	if len(baskets) > 1 {
		// must be fixed if happened
		return echo.ErrInternalServerError
	}

	if err = bh.repo.Delete(c.Request().Context(), basketrepo.GetCommand{ID: &id, UserID: &userID}); err != nil {
		return echo.ErrInternalServerError
	}

	// get all baskets again
	baskets, jsonbs := bh.repo.Get(c.Request().Context(), basketrepo.GetCommand{UserID: &userID})
	resultMap := make(map[string]interface{})
	resultMap["baskets"] = baskets
	resultMap["jsonbs"] = jsonbs
	return c.JSONPretty(http.StatusOK, resultMap, "  ")

}

func (bh *BasketHandler) Register(g *echo.Group) {
	g.GET("", bh.Get)
	g.POST("", bh.Create)
	g.PATCH(":id", bh.Update)
	g.GET(":id", bh.GetByID)
	g.DELETE(":id", bh.Delete)
}
