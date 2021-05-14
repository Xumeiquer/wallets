package api

import (
	"net/http"

	"encoding/json"

	"github.com/Xumeiquer/wallets/models"
	"github.com/Xumeiquer/wallets/server/middleware/persistence"
	"github.com/labstack/echo"
)

// UpdateWallet updates a wallet.
// Params:
//  - `:id`: Wallet ID which will be updated
//  - A JSON object which represents the Wallet model
// Returns:
//  A JSON object that represents the status of the operation.
//  This JSON object can be either from error message or successful response.
func UpdateWallet(c echo.Context) error {
	// Get the wallet ID from the URI
	walletID := c.Param("id")

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

	// Update wallet
	err := db.UpdateWallet(walletID, wallet)
	if err != nil {
		// Updating error
		errMsg := models.ErrorMessage{
			Errno: models.ERR_UPDATE_WALLET,
			Msg:   models.ERR_MESSAGES[models.ERR_UPDATE_WALLET],
		}
		errMsgByte, _ := json.Marshal(&errMsg)
		apiResponse := models.APIResponse{
			Type: "error",
			Msg:  errMsgByte,
		}
		return c.JSON(http.StatusOK, apiResponse)
	}

	// All good
	seccessMsg := models.ErrorMessage{
		Errno: models.MSG_STATUS_SUCCESS,
		Msg:   models.ERR_MESSAGES[models.MSG_STATUS_SUCCESS],
	}
	seccessMsgByte, _ := json.Marshal(&seccessMsg)
	apiResponse := models.APIResponse{
		Type: "success",
		Msg:  seccessMsgByte,
	}
	return c.JSON(http.StatusOK, apiResponse)
}

// func UpdateWalletWInsi(c echo.Context) error {
// 	walletID := c.Param("id")
// 	fundINSI := c.Param("insi")

// 	db, ok := c.Get("db").(*persistence.Database)
// 	if !ok {
// 		errMsg := models.ErrorMessage{
// 			Errno: models.ERR_DB_NOT_FOUND,
// 			Msg:   models.ERR_MESSAGES[models.ERR_DB_NOT_FOUND],
// 		}
// 		return c.JSON(http.StatusOK, errMsg)
// 	}

// 	err := db.UpdateWalletFunds(walletID, fundINSI)
// 	if err != nil {
// 		errMsg := models.ErrorMessage{
// 			Errno: models.ERR_UPDATE_WALLET,
// 			Msg:   models.ERR_MESSAGES[models.ERR_UPDATE_WALLET],
// 		}
// 		return c.JSON(http.StatusOK, errMsg)
// 	}

// 	errMsg := models.ErrorMessage{
// 		Errno: models.MSG_STATUS_SUCCESS,
// 		Msg:   models.ERR_MESSAGES[models.MSG_STATUS_SUCCESS],
// 	}
// 	return c.JSON(http.StatusOK, errMsg)
// }
