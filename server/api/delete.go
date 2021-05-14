package api

import (
	"net/http"

	"encoding/json"

	"github.com/Xumeiquer/wallets/models"
	"github.com/Xumeiquer/wallets/server/middleware/persistence"
	"github.com/labstack/echo"
)

// DeleteWallet deletes a wallet.
// Params:
//  - `:id`: is the Wallet ID
// Returns:
//  A JSON object that represents the status of the operation.
//  This JSON object can be either from error message or successful response.
func DeleteWallet(c echo.Context) error {
	// Get the wallet ID from the URI
	walletID := c.Param("id")
	// wallet := models.Wallet{}
	// if err := c.Bind(&wallet); err != nil {
	// 	errMsg := models.ErrorMessage{
	// 		Errno: models.ERR_NOT_A_WALLET,
	// 		Msg:   models.ERR_MESSAGES[models.ERR_NOT_A_WALLET],
	// 	}
	// 	return c.JSON(http.StatusOK, errMsg)
	// }

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

	// Delete the wallet based on the wallet ID
	err := db.DeleteWallet(walletID)
	if err != nil {
		// Error while deleting the wallet
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
	successMsg := models.ErrorMessage{
		Errno: models.MSG_STATUS_SUCCESS,
		Msg:   models.ERR_MESSAGES[models.MSG_STATUS_SUCCESS],
	}
	successMsgByte, _ := json.Marshal(&successMsg)
	apiResponse := models.APIResponse{
		Type: "success",
		Msg:  successMsgByte,
	}
	return c.JSON(http.StatusOK, apiResponse)
}
