package middleware

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Xumeiquer/wallets/models"
)

type APISetter interface {
	APISet(*API)
}

type APIRef struct{ *API }

func (cr *APIRef) APISet(c *API) { cr.API = c }

// API helps to communicate with the backend
type API struct {
	Endpoint string
	routes   map[string]string
}

const (
	R_READ_WALLETS  = "ReadWallets"
	R_CREATE_WALLET = "CreateWallet"
	R_READ_WALLET   = "ReadWallet"
	R_UPDATE_WALLET = "UpdateWallet"
	R_DELETE_WALLET = "DeleteWallet"
	R_READ_FUND     = "ReadFund"
)

func NewAPI(endpoint string) *API {
	return &API{
		Endpoint: endpoint,
		routes: map[string]string{
			R_READ_WALLETS:  fmt.Sprintf("http://%s/api/v1/wallet/", endpoint),
			R_CREATE_WALLET: fmt.Sprintf("http://%s/api/v1/wallet/", endpoint),
			R_READ_WALLET:   fmt.Sprintf("http://%s/api/v1/wallet/%%s", endpoint),
			R_UPDATE_WALLET: fmt.Sprintf("http://%s/api/v1/wallet/%%s", endpoint),
			R_DELETE_WALLET: fmt.Sprintf("http://%s/api/v1/wallet/%%s", endpoint),
			R_READ_FUND:     fmt.Sprintf("http://%s/api/v1/fund/%%s", endpoint),
		},
	}
}

// GetWallets get all wallets
func (api *API) GetWallets() ([]models.Wallet, error) {
	resp, err := http.Get(api.routes[R_READ_WALLETS])
	if err != nil {
		fmt.Printf("[API] ERROR - GET :: %s\n", err.Error())
		return []models.Wallet{}, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("[API] ERROR - READ_ALL_BODY :: %s\n", err.Error())
		return []models.Wallet{}, err
	}

	var apiResponse models.APIResponse
	err = json.Unmarshal(data, &apiResponse)
	if err != nil {
		fmt.Printf("[API] ERROR - JSON_UNMARSHAL :: API_RESPONSE :: %s\n", err.Error())
		return []models.Wallet{}, err
	}

	switch apiResponse.Type {

	case "error":
		var errMsg models.ErrorMessage
		err = json.Unmarshal(apiResponse.Msg, &errMsg)
		if err != nil {
			fmt.Printf("[API] ERROR - JSON_UNMARSHAL :: ERROR_MESSAGE :: %s\n", err.Error())
		} else {
			fmt.Printf("[API] INFO - JSON_UNMARSHAL :: %v\n", errMsg)
		}
		return []models.Wallet{}, errors.New(errMsg.Msg)

	case "wallets":
		var wallets []models.Wallet
		err = json.Unmarshal(apiResponse.Msg, &wallets)
		if err != nil {
			fmt.Printf("[API] ERROR - JSON_UNMARSHAL :: WALLETS :: %s\n", err.Error())
			return []models.Wallet{}, err
		}
		return wallets, err
	}

	return []models.Wallet{}, nil
}

func (api *API) CreateWallet(walletName string) (models.Wallet, error) {
	newWallet, err := models.NewWallet(walletName)
	if err != nil {
		fmt.Printf("[API] ERROR - NEW_WALLET :: %s\n", err.Error())
		return models.Wallet{}, err
	}

	byteData, err := json.Marshal(&newWallet)
	if err != nil {
		fmt.Printf("[API] ERROR - JSON_MARSHAL :: %s\n", err.Error())
		return models.Wallet{}, err
	}

	data := bytes.NewBuffer(byteData)
	resp, err := http.Post(api.routes[R_CREATE_WALLET], "application/json", data)
	if err != nil {
		fmt.Printf("[API] ERROR - POST :: %s\n", err.Error())
		return models.Wallet{}, err
	}

	defer resp.Body.Close()
	byteData, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("[API] ERROR - READ_ALL_BODY :: %s\n", err.Error())
		return models.Wallet{}, err
	}

	var apiResponse models.APIResponse
	err = json.Unmarshal(byteData, &apiResponse)
	if err != nil {
		fmt.Printf("[API] ERROR - JSON_UNMARSHAL :: API_RESPONSE :: %s\n", err.Error())
		return models.Wallet{}, err
	}

	switch apiResponse.Type {

	case "error":
		var errMsg models.ErrorMessage
		err = json.Unmarshal(apiResponse.Msg, &errMsg)
		if err != nil {
			fmt.Printf("[API] ERROR - JSON_UNMARSHAL :: ERROR_MESSAGE :: %s\n", err.Error())
		} else {
			fmt.Printf("[API] INFO - JSON_UNMARSHAL :: %v\n", errMsg)
		}
		return models.Wallet{}, err

	case "wallet":
		var wallets models.Wallet
		err = json.Unmarshal(apiResponse.Msg, &wallets)
		if err != nil {
			fmt.Printf("[API] ERROR - JSON_UNMARSHAL :: WALLET :: %s\n", err.Error())
			return models.Wallet{}, err
		}
		return wallets, err
	}

	return models.Wallet{}, nil
}

func (api *API) UpdateWallet(wallet models.Wallet) error {
	byteData, err := json.Marshal(&wallet)
	if err != nil {
		fmt.Printf("[API] ERROR - JSON_MARSHAL :: %s\n", err.Error())
		return err
	}

	url := fmt.Sprintf(api.routes[R_UPDATE_WALLET], wallet.GetID())
	data := bytes.NewBuffer(byteData)

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, url, data)
	if err != nil {
		fmt.Printf("[API] ERROR - PUT :: NEW_REQUEST :: %s\n", err.Error())
		return err
	}

	req.Header.Set("content-type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("[API] ERROR - PUT :: %s\n", err.Error())
		return err
	}
	defer resp.Body.Close()
	byteData, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("[API] ERROR - READ_ALL_BODY :: %s\n", err.Error())
		return err
	}

	var apiResponse models.APIResponse
	err = json.Unmarshal(byteData, &apiResponse)
	if err != nil {
		fmt.Printf("[API] ERROR - JSON_UNMARSHAL :: API_RESPONSE :: %s\n", err.Error())
		return err
	}

	switch apiResponse.Type {

	case "error":
		var errMsg models.ErrorMessage
		err = json.Unmarshal(apiResponse.Msg, &errMsg)
		if err != nil {
			fmt.Printf("[API] ERROR - JSON_UNMARSHAL :: ERROR_MESSAGE :: %s\n", err.Error())
		} else {
			fmt.Printf("[API] WARN - JSON_UNMARSHAL :: %v\n", errMsg)
		}
		return err

	case "success":
		var errMsg models.ErrorMessage
		err = json.Unmarshal(apiResponse.Msg, &errMsg)
		if err != nil {
			fmt.Printf("[API] ERROR - JSON_UNMARSHAL :: ERROR_MESSAGE :: %s\n", err.Error())
		} else {
			fmt.Printf("[API] INFO - JSON_UNMARSHAL :: %v\n", errMsg)
		}
		return nil
	}

	return nil
}

func (api *API) DeleteWallet(wallet models.Wallet) error {
	url := fmt.Sprintf(api.routes[R_DELETE_WALLET], wallet.GetID())

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		fmt.Printf("[API] ERROR - NEW_REQUEST :: %s\n", err.Error())
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("[API] ERROR - DELETE :: %s\n", err.Error())
		return err
	}
	defer resp.Body.Close()
	byteData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("[API] ERROR - REAL_ALL_BODY :: %s\n", err.Error())
		return err
	}

	var apiResponse models.APIResponse
	err = json.Unmarshal(byteData, &apiResponse)
	if err != nil {
		fmt.Printf("[API] ERROR - JSON_UNMARSHAL :: API_RESPONSE :: %s\n", err.Error())
		return err
	}

	switch apiResponse.Type {

	case "error":
		var errMsg models.ErrorMessage
		err = json.Unmarshal(apiResponse.Msg, &errMsg)
		if err != nil {
			fmt.Printf("[API] ERROR - JSON_UNMARSHAL :: ERROR_MESSAGE :: %s\n", err.Error())
		} else {
			fmt.Printf("[API] WARN - JSON_UNMARSHAL :: %v\n", errMsg)
		}
		return err

	case "success":
		var errMsg models.ErrorMessage
		err = json.Unmarshal(apiResponse.Msg, &errMsg)
		if err != nil {
			fmt.Printf("[API] ERROR - JSON_UNMARSHAL :: ERROR_MESSAGE :: %s\n", err.Error())
		} else {
			fmt.Printf("[API] WARN - JSON_UNMARSHAL :: %v\n", errMsg)
		}
		return nil
	}

	return nil
}

func (api *API) GetFund(insi string) (models.FundData, error) {
	resp, err := http.Get(fmt.Sprintf(api.routes[R_READ_FUND], insi))
	if err != nil {
		fmt.Printf("[API] ERROR - GET :: %s\n", err.Error())
		return models.FundData{}, err
	}
	defer resp.Body.Close()

	var byteData []byte
	byteData, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("[API] ERROR - READ_ALL_BODY :: %s\n", err.Error())
		return models.FundData{}, err
	}

	var apiResponse models.APIResponse
	err = json.Unmarshal(byteData, &apiResponse)
	if err != nil {
		fmt.Printf("[API] ERROR - JSON_UNMARSHAL :: API_RESPONSE :: %s\n", err.Error())
		return models.FundData{}, err
	}

	switch apiResponse.Type {

	case "error":
		var errMsg models.ErrorMessage
		err = json.Unmarshal(apiResponse.Msg, &errMsg)
		if err != nil {
			fmt.Printf("[API] ERROR - JSON_UNMARSHAL :: ERROR_MESSAGE :: %s\n", err.Error())
		} else {
			fmt.Printf("[API] WARN - JSON_UNMARSHAL :: %v\n", errMsg)
		}
		return models.FundData{}, err

	case "fund":
		var fund models.FundData
		err = json.Unmarshal(apiResponse.Msg, &fund)
		if err != nil {
			fmt.Printf("[API] ERROR - JSON_UNMARSHAL :: FUND_DATA :: %s\n", err.Error())
			return models.FundData{}, err
		}
		return fund, err
	}
	return models.FundData{}, nil
}
