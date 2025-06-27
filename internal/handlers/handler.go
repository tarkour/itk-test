package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/tarkour/itk-test/pkg/database"
	_ "github.com/tarkour/itk-test/pkg/database"
)

func HandlerGet(ctx echo.Context, db *database.DBConn, walletid int) error {

	start := time.Now()
	defer func() {
		log.Printf("HandlerGet execution time: %v", time.Since(start))
	}()

	data, err := db.DBGet(walletid)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "database error: "+err.Error())
	}

	res := strconv.Itoa(data)

	return ctx.String(http.StatusOK, res)

}

func HandlerPost(ctx echo.Context, db *database.DBConn) error {
	start := time.Now()
	defer func() {
		log.Printf("HandlerGet execution time: %v", time.Since(start))
	}()

	var wallet database.Wallet
	if err := ctx.Bind(&wallet); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	if wallet.WalletID == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "WalletID is required")
	}

	if wallet.Amount <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Amount must be positive")
	}

	if err := db.DBPost(&wallet); err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}

	return ctx.NoContent(http.StatusOK)
}
