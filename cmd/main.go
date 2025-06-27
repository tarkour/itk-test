package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tarkour/itk-test/internal/handlers"
	"github.com/tarkour/itk-test/migrations"
	"github.com/tarkour/itk-test/pkg/database"
)

func main() {

	dbConn, err := database.ConnectDB("./")
	if err != nil {
		log.Fatalf("db connection error: %v", err)
	}

	defer dbConn.Close()

	migrations.Migrate()
	fmt.Println("migrations are done successfully")

	fmt.Println("Database connected successfully")

	serverAddr := os.Getenv("SERVER_ADDRESS")
	if serverAddr == "" {
		serverAddr = ":8080"
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet},
	}))

	e.GET("/", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "server is running!")
	})
	e.GET("/test", Handler)
	e.GET("api/v1/wallets/:walletid", func(ctx echo.Context) error {
		walletID, err := strconv.Atoi(ctx.Param("walletid"))
		if err != nil {
			return echo.NewHTTPError(400, "Invalid wallet ID format")
		}
		return handlers.HandlerGet(ctx, dbConn, walletID)
	})
	e.POST("/api/v1/wallet", func(ctx echo.Context) error {
		return handlers.HandlerPost(ctx, dbConn)
	})

	if err := e.Start(serverAddr); err != nil {
		e.Logger.Fatal("Failed to start server: ", err)
	}

}

func Handler(ctx echo.Context) error {
	err := ctx.String(http.StatusOK, "test")
	if err != nil {
		return err
	}

	return nil
}
