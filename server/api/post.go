package api

import (
	"net/http"

	"encoding/json"

	"github.com/Xumeiquer/wallets/models"
	"github.com/Xumeiquer/wallets/server/middleware/persistence"
	"github.com/labstack/echo"
)

// CreateWallet creates a new wallet.
// Params:
//  - A JSON object which represents the Wallet model
// Returns:
//  A JSON object that represents the status of the operation.
//  This JSON object can be either from error message or successful response.
func CreateWallet(c echo.Context) error {
	// Bind PUT payload as wallet
	wallet := models.Wallet{}
	if err := c.Bind(&wallet); err != nil {
		// The payload is not a wallet
		errMsg := models.ErrorMessage{
			Errno: models.ERR_NOT_A_WALLET,
			Msg:   models.ERR_MESSAGES[models.ERR_NOT_A_WALLET],
		}
		errMsgByte, _ := json.Marshal(&errMsg)
		apiResponse := models.APIResponse{
			Type: "error",
			Msg:  errMsgByte,
		}
		return c.JSON(http.StatusOK, apiResponse)
	}

	// Grab the DB from the context
	db, ok := c.Get("db").(*persistence.Database)
	if !ok {
		// Weird. The DB is not avaliable
		errMsg := models.ErrorMessage{
			Errno: models.ERR_DB_NOT_FOUND,
			Msg:   models.ERR_MESSAGES[models.ERR_DB_NOT_FOUND],
		}
		errMsgByte, _ := json.Marshal(&errMsg)
		apiResponse := models.APIResponse{
			Type: "error",
			Msg:  errMsgByte,
		}
		return c.JSON(http.StatusOK, apiResponse)
	}

	// TODO: Perform extra valdations

	// w, err := models.NewWallet(wallet.GetName())
	// if err != nil {
	// 	errMsg := models.ErrorMessage{
	// 		Errno: models.ERR_NEW_WALLET,
	// 		Msg:   models.ERR_MESSAGES[models.ERR_NEW_WALLET],
	// 	}
	// 	return c.JSON(http.StatusOK, errMsg)
	// }

	err := db.UpdateWallet(wallet.GetID(), wallet)
	if err != nil {
		// Error saving the wallet
		errMsg := models.ErrorMessage{
			Errno: models.ERR_INSERT_WALLET,
			Msg:   models.ERR_MESSAGES[models.ERR_INSERT_WALLET],
		}
		errMsgByte, _ := json.Marshal(&errMsg)
		apiResponse := models.APIResponse{
			Type: "error",
			Msg:  errMsgByte,
		}
		return c.JSON(http.StatusOK, apiResponse)
	}

	// All good
	walletByte, _ := json.Marshal(&wallet)
	apiResponse := models.APIResponse{
		Type: "wallet",
		Msg:  walletByte,
	}
	return c.JSON(http.StatusOK, apiResponse)
}
