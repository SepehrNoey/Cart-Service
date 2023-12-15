package handler

import (
	"net/http"

	"github.com/SepehrNoey/Cart-Service/internal/domain/repository/basketrepo"
	"github.com/labstack/echo/v4"
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

}

func (bh *BasketHandler) Register(g *echo.Group) {
	g.GET("", bh.Get)
	g.POST("", bh.Create)
	g.PATCH(":id", bh.Update)
	g.GET(":id", bh.GetByID)
	g.DELETE(":id", bh.Delete)
}
