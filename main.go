package main

import (
	"fmt"
	"log"
	"time"

	"github.com/SepehrNoey/Cart-Service/internal/domain/repository/basketrepo"
	"github.com/SepehrNoey/Cart-Service/internal/domain/repository/userrepo"
	"github.com/SepehrNoey/Cart-Service/internal/infra/http/auth"
	"github.com/SepehrNoey/Cart-Service/internal/infra/http/handler"
	"github.com/SepehrNoey/Cart-Service/internal/infra/repository/basket/basketsql"
	"github.com/SepehrNoey/Cart-Service/internal/infra/repository/user/usersql"
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

	if err = db.AutoMigrate(&basketsql.BasketDTO{}, &usersql.UserDTO{}); err != nil {
		fmt.Printf("failed to automigrate: %v", err)
	}

	app := echo.New()
	var basketRepo basketrepo.Repository = basketsql.New(db)
	var userRepo userrepo.Repository = usersql.New(db)

	bh := handler.NewBasketHandler(basketRepo)
	bh.Register(app.Group("basket/"))
	uh := handler.NewUserHandler(userRepo)
	uh.Register(app.Group("/"))

	secretKey := []byte("secret-key-of-basket-service-for-jwt-authentication")
	expDur := time.Minute * 5
	auth.SetJWTConfig(secretKey, expDur)

	if err := app.Start("0.0.0.0:2023"); err != nil {
		log.Fatalf("server failed to start %v", err)
	}
}
