package main

import (
	"fmt"
	"log"

	"github.com/SepehrNoey/Cart-Service/internal/domain/repository/basketrepo"
	"github.com/SepehrNoey/Cart-Service/internal/infra/http/handler"
	"github.com/SepehrNoey/Cart-Service/internal/infra/repository/basketsql"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=postgres password=0150188511 dbname=basket-service-go port=5432 sslmode=disable TimeZone=Asia/Tehran"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("error connecting to database: %v", err)
	}

	if err = db.AutoMigrate(new(basketsql.BasketDTO)); err != nil {
		fmt.Printf("failed to automigrate: %v", err)
	}

	app := echo.New()
	var repo basketrepo.Repository = basketsql.New(db)

	h := handler.NewBasketHandler(repo)
	h.Register(app.Group("basket/"))

	if err := app.Start("0.0.0.0:2023"); err != nil {
		log.Fatalf("server failed to start %v", err)
	}
}
