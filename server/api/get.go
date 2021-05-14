package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"encoding/json"

	"github.com/Xumeiquer/wallets/models"
	"github.com/Xumeiquer/wallets/server/middleware/persistence"
	"github.com/antchfx/htmlquery"
	"github.com/labstack/echo"
	"golang.org/x/net/html"
)

// ReadWallet reads a wallet. When the `id` is defined in the request URI
// this method returns the wallet with that `id`, otherwise ReadWallet
// returns all wallets stored in the database.
// Params:
//  - `:id`: Wallet ID
// Returns:
//  A JSON object that represents the status of the operation.
//  This JSON object can be either from error message or successful response.
func ReadWallet(c echo.Context) error {
	// Get the wallet ID from the URI
	walletID := c.Param("id")

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

	// Get wallet based on the wallet ID
	wallet, err := db.GetWallet(walletID)
	if err != nil {
		// Error geting the wallet
		errMsg := models.ErrorMessage{
			Errno: models.ERR_WALLET_NOT_FOUND,
			Msg:   models.ERR_MESSAGES[models.ERR_WALLET_NOT_FOUND],
		}
		errMsgByte, _ := json.Marshal(&errMsg)
		apiResponse := models.APIResponse{
			Type: "error",
			Msg:  errMsgByte,
		}
		return c.JSON(http.StatusOK, apiResponse)
	}

	// There are no error, but there are no wallets to send over
	if len(wallet) == 0 {
		// Is this needed?
		errMsg := models.ErrorMessage{
			Errno: models.ERR_WALLET_NOT_FOUND,
			Msg:   models.ERR_MESSAGES[models.ERR_WALLET_NOT_FOUND],
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
		Type: "wallets",
		Msg:  walletByte,
	}
	return c.JSON(http.StatusOK, apiResponse)
}

func ReadFund(c echo.Context) error {
	var (
		docFT, docMO, node *html.Node
		aux                float64
		err                error
	)

	insi := c.Param("insi")

	docFT, err = htmlquery.LoadURL(fmt.Sprintf("https://markets.ft.com/data/funds/tearsheet/performance?s=%s:EUR", insi))
	if err != nil {
		fmt.Printf("[API] ERROR - GET :: %s\n", err.Error())
		errMsg := models.ErrorMessage{
			Errno: models.ERR_INSI_NOT_FOUND,
			Msg:   models.ERR_MESSAGES[models.ERR_INSI_NOT_FOUND],
		}
		errMsgByte, _ := json.Marshal(&errMsg)
		apiResponse := models.APIResponse{
			Type: "error",
			Msg:  errMsgByte,
		}
		return c.JSON(http.StatusOK, apiResponse)
	}

	docMO, err = htmlquery.LoadURL(fmt.Sprintf("http://performance.morningstar.com/Performance/fund/trailing-total-returns.action?t=%s&s=%s", insi, insi))
	if err != nil {
		fmt.Printf("[API] ERROR - GET :: %s\n", err.Error())
		errMsg := models.ErrorMessage{
			Errno: models.ERR_INSI_NOT_FOUND,
			Msg:   models.ERR_MESSAGES[models.ERR_INSI_NOT_FOUND],
		}
		errMsgByte, _ := json.Marshal(&errMsg)
		apiResponse := models.APIResponse{
			Type: "error",
			Msg:  errMsgByte,
		}
		return c.JSON(http.StatusOK, apiResponse)
	}

	fund := models.FundData{
		INSI: insi,
	}

	node = htmlquery.FindOne(docFT, "//div[@class='mod-tearsheet-overview__header']/h1[1]/text()")
	fund.Name = node.Data

	node = htmlquery.FindOne(docFT, "//div[@class='mod-tearsheet-overview__quote']/ul/li[1]/span[2]/text()")
	aux, err = strconv.ParseFloat(node.Data, 32)
	if err != nil {
		fund.Nav = 0
	} else {
		fund.Nav = float32(aux)
	}

	node = htmlquery.FindOne(docFT, "//ul[@class='mod-tearsheet-overview__quote__bar']/li[2]/span[@class='mod-ui-data-list__value']/span[1]/text()")
	values := strings.Split(node.Data, "/")
	for idx, value := range values {
		aux, err = strconv.ParseFloat(value, 32)
		if err != nil {
			fund.DayChange[idx] = 0
		} else {
			fund.DayChange[idx] = float32(aux)
		}
	}

	node = htmlquery.FindOne(docFT, "//div[@class='mod-tearsheet-overview__quote']/div[1]/text()")
	values = strings.Split(node.Data, ",")
	if len(values) == 2 {
		aux := strings.Replace(values[1], "as of", "", -1)
		aux = strings.Replace(aux, ".", "", -1)
		aux = strings.TrimSpace(aux)
		fund.TimeDelayed = aux
	} else {
		fund.TimeDelayed = ""
	}

	node = htmlquery.FindOne(docFT, "//div[@class='mod-ui-table--freeze-pane__scroll-container']/table/tbody/tr[1]/td[7]/span[1]/text()")
	fund.OneMonthS = node.Data

	node = htmlquery.FindOne(docFT, "//div[@class='mod-ui-table--freeze-pane__scroll-container']/table/tbody/tr[1]/td[6]/span[1]/text()")
	fund.ThreeMonthS = node.Data

	node = htmlquery.FindOne(docFT, "//div[@class='mod-ui-table--freeze-pane__scroll-container']/table/tbody/tr[1]/td[5]/span[1]/text()")
	fund.SixMonthS = node.Data

	node = htmlquery.FindOne(docFT, "//div[@class='mod-ui-table--freeze-pane__scroll-container']/table/tbody/tr[1]/td[4]/span[1]/text()")
	fund.OneYearS = node.Data

	node = htmlquery.FindOne(docMO, "//th[@class='col_head_lbl']/span[1]/text()")
	value := strings.TrimSpace(node.Data)
	value = strings.Replace(value, "(", "", -1)
	value = strings.Replace(value, ")", "", -1)
	value = strings.Replace(value, "&nbsp;", "", -1)
	fund.Date, _ = time.Parse("01/02/2006", value)

	node = htmlquery.FindOne(docMO, "/html/body/table/tbody[1]/tr[1]/td[1]/text()")
	aux, err = strconv.ParseFloat(node.Data, 32)
	if err != nil {
		fund.OneDay = 0
	} else {
		fund.OneDay = float32(aux)
	}

	node = htmlquery.FindOne(docMO, "/html/body/table/tbody[1]/tr[1]/td[2]/text()")
	aux, err = strconv.ParseFloat(node.Data, 32)
	if err != nil {
		fund.OneWeek = 0
	} else {
		fund.OneWeek = float32(aux)
	}

	node = htmlquery.FindOne(docMO, "/html/body/table/tbody[1]/tr[1]/td[3]/text()")
	aux, err = strconv.ParseFloat(node.Data, 32)
	if err != nil {
		fund.OneMonth = 0
	} else {
		fund.OneMonth = float32(aux)
	}

	node = htmlquery.FindOne(docMO, "/html/body/table/tbody[1]/tr[1]/td[4]/text()")
	aux, err = strconv.ParseFloat(node.Data, 32)
	if err != nil {
		fund.ThreeMonth = 0
	} else {
		fund.ThreeMonth = float32(aux)
	}

	node = htmlquery.FindOne(docMO, "/html/body/table/tbody[1]/tr[1]/td[6]/text()")
	aux, err = strconv.ParseFloat(node.Data, 32)
	if err != nil {
		fund.OneYear = 0
	} else {
		fund.OneYear = float32(aux)
	}

	node = htmlquery.FindOne(docMO, "/html/body/table/tbody[1]/tr[1]/td[7]/text()")
	aux, err = strconv.ParseFloat(node.Data, 32)
	if err != nil {
		fund.ThreeYears = 0
	} else {
		fund.ThreeYears = float32(aux)
	}

	node = htmlquery.FindOne(docMO, "/html/body/table/tbody[1]/tr[1]/td[8]/text()")
	aux, err = strconv.ParseFloat(node.Data, 32)
	if err != nil {
		fund.FiveYears = 0
	} else {
		fund.FiveYears = float32(aux)
	}

	fundByte, _ := json.Marshal(&fund)
	apiResponse := models.APIResponse{
		Type: "fund",
		Msg:  fundByte,
	}
	return c.JSON(http.StatusOK, apiResponse)
}
