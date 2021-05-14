package router

import (
	"github.com/Xumeiquer/wallets/server/api"
	"github.com/labstack/echo"
)

const (
	APIPREFIX = "/api"
	APIV1     = "/v1"
	APILATEST = APIV1
)

// SetRouting defines the API routing
func SetRouting(e *echo.Echo) {
	v1 := e.Group(APIPREFIX).Group(APIV1)
	{
		wallet := v1.Group("/wallet")
		{
			// Wallet CRUD
			wallet.GET("/", api.ReadWallet)
			wallet.POST("/", api.CreateWallet)
			wallet.GET("/:id", api.ReadWallet)
			wallet.PUT("/:id", api.UpdateWallet)
			wallet.DELETE("/:id", api.DeleteWallet)

			// wallet.PUT("/:id/fund/:insi", api.UpdateWalletWInsi)
		}

		fund := v1.Group("/fund")
		{
			fund.GET("/:insi", api.ReadFund)
		}
	}
}
