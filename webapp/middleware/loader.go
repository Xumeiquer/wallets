package middleware

import (
	"errors"
	"fmt"

	"github.com/Xumeiquer/wallets/models"
)

type LoaderSetter interface {
	LoaderSet(*Loader)
}

type LoaderRef struct{ *Loader }

func (cr *LoaderRef) LoaderSet(c *Loader) { cr.Loader = c }

// Loader helps to communicate with the backend
type Loader struct {
	S *State
	A *API
}

// NewLoader returns a new loader. This needs the API and State
func NewLoader(st *State, api *API) *Loader {
	return &Loader{
		S: st,
		A: api,
	}
}

// UpdateFundsData update all FundData for the whole wallet
func (lo *Loader) UpdateFundsData(activeWallet *models.Wallet) error {
	for insi := range activeWallet.GetFunds() {
		err := lo.UpdateFundData(activeWallet, insi)
		if err != nil {
			return err
		}
	}
	return nil

}

func (lo *Loader) UpdateFundData(activeWallet *models.Wallet, insi string) error {
	// TODO: This may need to implement a force option to avoid using the cache
	var (
		fund models.FundData
		err  error
		ok   bool
	)

	// Fund needs to be added to wallet
	needToAdd := true

	// Search whether the fund data is in the cache
	fundData, err := lo.S.Get(insi)
	if err == nil {
		// Found it
		fund, ok = fundData.(models.FundData)
		if !ok {
			fmt.Printf("[WALLET] ERROR :: Casting\n")
			return errors.New("data not valid for fund model")
		}
		// Fund was already loaded
		needToAdd = false
	} else {
		// Do a request to get the latest
		fund, err = lo.A.GetFund(insi)
		if err != nil {
			fmt.Printf("[WALLET] ERROR :: SetState %s\n", err.Error())
			return err
		}
	}

	// In case the INSI is new it should be added to the wallet
	if needToAdd {
		activeWallet.UpdateFundData(fund)

		err = lo.S.Set(insi, fund)
		if err != nil {
			fmt.Printf("[WALLET] ERROR :: SetState %s\n", err.Error())
			return err
		}
	}

	return nil
}
