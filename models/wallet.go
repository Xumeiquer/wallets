package models

import (
	"errors"
	"strings"
	"time"

	"github.com/Xumeiquer/wallets/server/utils"
)

// Wallet represents a wallet
type Wallet struct {
	ID        string
	Name      string
	Funds     Funds
	Ops       Operations
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Operations represents the operations that were taken by a Fund
type Operations map[string][]Operation

// Funds contains the dymanic data for each Fund
type Funds map[string]FundData

// Operation represents one operation
type Operation struct {
	ID       string
	FundInsi string
	Date     time.Time
	Amount   float32
	NAV      float32
	Shares   float32
}

// NewWallet creates a new wallet with a defined name
func NewWallet(name string) (Wallet, error) {
	if !isValidName(name) {
		return Wallet{}, errors.New("name not valid")
	}

	return Wallet{
		ID:        utils.NewULID(),
		Name:      name,
		Funds:     map[string]FundData{},
		Ops:       make(Operations),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (w *Wallet) update() error {
	w.UpdatedAt = time.Now()
	return nil
}

// GetCopy creates a copy of the wallet
func (w Wallet) GetCopy() Wallet {
	return Wallet{
		ID:        w.ID,
		Name:      w.Name,
		Funds:     CopyFundsMap(w.Funds),
		Ops:       CopyOperationsMap(w.Ops),
		CreatedAt: w.CreatedAt,
		UpdatedAt: w.UpdatedAt,
	}
}

// GetID returns the wallet ID
func (w Wallet) GetID() string {
	return w.ID
}

// GetName returns the wallet name
func (w Wallet) GetName() string {
	return w.Name
}

// SetName sets the wallet name
func (w *Wallet) SetName(name string) error {
	if !isValidName(name) {
		return errors.New("name not valid")
	}

	w.Name = name
	return w.update()
}

// GetFunds returns the Funds associated to the wallet
func (w Wallet) GetFunds() Funds {
	return w.Funds
}

// GetFundData returns the dynamic data associated to each Fund
func (w Wallet) GetFundData(insi string) FundData {
	return w.Funds[insi]
}

// GetFundNameByInsi retruns the Fund name for a given INSI
func (w Wallet) GetFundNameByInsi(insi string) string {
	fund, ok := w.Funds[insi]
	if ok {
		return fund.GetName()
	}
	return ""
}

// GetFundInsiByName returns the Fund INSI for a given name
func (w Wallet) GetFundInsiByName(name string) string {
	for _, fundData := range w.Funds {
		if fundData.GetName() == name {
			return fundData.GetINSI()
		}
	}

	return ""
}

// RemoveFundByInsi removes a Fund from the wallet
func (w *Wallet) RemoveFundByInsi(insi string) {
	delete(w.Funds, insi)
	delete(w.Ops, insi)
}

// HasFund returns true if the Fund is associated to the wallet, otherwise it
// returns false
func (w Wallet) HasFund(insi string) bool {
	if !isValidName(insi) {
		return false
	}
	_, ok := w.Funds[insi]

	return ok
}

func (w *Wallet) UpdateOperation(opID string, newOp Operation) error {
	for _, ops := range w.Ops {
		for idx := 0; idx < len(ops); idx++ {
			if ops[idx].GetID() == opID {
				ops[idx].SetDate(newOp.GetDate())
				ops[idx].SetInsi(newOp.GetInsi())
				ops[idx].SetAmount(newOp.GetAmount())
				ops[idx].SetNav(newOp.GetNav())
				ops[idx].SetShares(newOp.GetShares())
				return nil
			}
		}
	}
	return errors.New("operation not found")
}

// RemoveOperation removed an operation from the wallet
func (w Wallet) RemoveOperation(opID string) error {
	for insi, ops := range w.Ops {
		for idx := 0; idx < len(ops); idx++ {
			if ops[idx].GetID() == opID {
				w.Ops[insi] = append(w.Ops[insi][:idx], w.Ops[insi][idx+1:]...)
				return nil
			}
		}
	}
	return errors.New("operation not found")
}

// AddFund adds the fund and its dynamic data to the wallet
func (w *Wallet) AddFund(fund FundData) error {
	if _, ok := w.Funds[fund.GetINSI()]; !ok {
		w.Funds[fund.GetINSI()] = fund
		return w.update()
	}

	return errors.New("fund already set")
}

// UpdateFundData updates the Fund dynamic data
func (w *Wallet) UpdateFundData(fund FundData) error {
	insi := fund.GetINSI()
	if _, ok := w.Funds[insi]; ok {
		w.Funds[insi] = fund
		return nil
	}

	return errors.New("fund not found")
}

// UpdateFundDataAndOps updates the Fund dynamic data and its operations
func (w *Wallet) UpdateFundDataAndOps(oldInsi string, fund FundData) error {
	if _, ok := w.Funds[oldInsi]; ok {
		w.Funds[oldInsi] = fund
		return nil
	}

	fundOps, _ := w.GetOperationsFor(oldInsi)
	for _, op := range fundOps {
		op.SetInsi(fund.GetINSI())
	}

	return errors.New("fund not found")
}

// GetOperations returns the operations
func (w Wallet) GetOperations() Operations {
	return w.Ops
}

// HasOperation returns true is the wallet has the operation, otherwise
// it returns false
func (w Wallet) HasOperation(op Operation) bool {
	for _, f := range w.Ops {
		for _, o := range f {
			if o.ID == op.ID {
				return true
			}
		}
	}
	return false
}

// GetOperationsFor returns all the operations associated to a Fund
func (w Wallet) GetOperationsFor(insi string) ([]Operation, error) {
	if !w.HasFund(insi) {
		return []Operation{}, errors.New("fund not found")
	}

	return w.Ops[insi], nil
}

// AddOperation adds an operation to the wallet
func (w *Wallet) AddOperation(op Operation) error {
	insi := op.GetInsi()
	if _, ok := w.Ops[insi]; ok {
		for _, o := range w.Ops[insi] {
			if o.ID == op.ID {
				return errors.New("Operation already set")
			}
		}
		w.Ops[insi] = append(w.Ops[insi], op)
	} else {
		w.Ops[insi] = []Operation{}
		w.Ops[insi] = append(w.Ops[insi], op)
	}

	return w.update()
}

// Capitalization returns the wallet capitalization. The return object is
// a map which each key is a Fund INSI and the value is slice if 3 elements.
// The first one represents the shares, the second the inputs, and the third one
// is the wallet's worth.
func (w Wallet) Capitalization() map[string][]float32 {
	res := map[string][]float32{}

	for insi, fundData := range w.Funds {
		ops := w.Ops[insi]
		if len(ops) <= 0 {
			// Fund without operations
			continue
		}

		shares := float32(0)
		inputs := float32(0)
		worth := float32(0)

		for _, op := range ops {
			// Add all the operation's shares
			shares += op.GetShares()
			// Add all the operation's inputs
			inputs += op.GetAmount()
		}

		// All the shares times its market price give us the Fund investment worth
		worth = shares * fundData.GetNav()

		res[insi] = []float32{shares, inputs, worth}
	}

	return res
}

// Profitability returns the wallet profitability. The return object is
// a map which each key is a Fund INSI and the value is slice if 4 elements.
// The first one is the profitability which is the difference between input
// and wotrh, the second one is the percentage for the profitability, the
// third one is the TWR (Time-Weighted Return). Finally, the fourth one is
// MWR Money-Weighted Return.
func (w Wallet) Profitability() map[string][]float32 {
	res := map[string][]float32{}

	for fundInsi, values := range w.Capitalization() {
		// Simple profitability
		res[fundInsi] = append(res[fundInsi], values[2]-values[1])
		// Percentage of profitability
		res[fundInsi] = append(res[fundInsi], (values[2]-values[1])/values[1])
	}

	for fundInsi, fundData := range w.Funds {
		ops, ok := w.Ops[fundInsi]
		if !ok {
			// Fund without operations
			continue
		}

		if len(ops) > 0 {
			firstOp := ops[0]
			// TWR
			res[fundInsi] = append(res[fundInsi], fundData.GetNav()/(firstOp.GetNav()-1))
			// MWR --> TODO
			res[fundInsi] = append(res[fundInsi], 0.0)
		}
	}

	return res
}

// Allocations returns the percentage of assets allocation in the wallet. The
// result is a map where the key is the Fund INSI and the value is a slice for
// the initial allocation and the actual allocation.
func (w Wallet) Allocations() map[string][]float32 {
	res := map[string][]float32{}

	if len(w.Ops) == 0 {
		// There are no operations at all
		return res
	}

	totalAllocation := float32(0)
	totalWorth := float32(0)

	for fundInsi, fundData := range w.Funds {
		fundOps, ok := w.Ops[fundInsi]
		if !ok {
			// Fund without operations
			continue
		}

		for _, op := range fundOps {
			// Counting the total input contributed
			totalAllocation += op.GetAmount()
			// Counting the wallet worth
			totalWorth += op.GetShares() * fundData.GetNav()
		}
	}

	for fundInsi, fundData := range w.Funds {
		fundOps, ok := w.Ops[fundInsi]
		if !ok {
			// Fund without operations
			continue
		}

		fundAllocation := float32(0)
		fundShares := float32(0)
		fundWorth := float32(0)

		for _, op := range fundOps {
			// Counting the input contributed by fund
			fundAllocation += op.GetAmount()
			// Counting the Fund shares
			fundShares += op.GetShares()
		}
		// Fund worth
		fundWorth = fundShares * fundData.GetNav()

		// Initial allocation
		res[fundInsi] = append(res[fundInsi], (fundAllocation/totalAllocation)*100)
		// Actual allocation
		res[fundInsi] = append(res[fundInsi], (fundWorth/totalWorth)*100)
	}

	return res
}

// GetAssetsAllocationPOJO retuns an object which can be consumed by Charts.js.
// The object is usefull for plotting the allocations
func (w Wallet) GetAssetsAllocationPOJO() map[string]interface{} {
	labels := []interface{}{}
	data := []interface{}{}
	colors := []interface{}{}

	for _, c := range chartColors {
		colors = append(colors, c)
	}

	for fundInsi, allocation := range w.Allocations() {
		labels = append(labels, w.GetFundNameByInsi(fundInsi))
		data = append(data, allocation[1])
	}

	return map[string]interface{}{
		"data": map[string]interface{}{
			"labels": labels,
			"datasets": []interface{}{
				map[string]interface{}{
					"label":           "Funds",
					"data":            data,
					"backgroundColor": colors,
				},
			},
		},
	}
}

// GetHistoricWalletPOJO retuns an object which can be consumed by Charts.js.
// The object is usefull for plotting the inputs history
func (w Wallet) GetHistoricWalletPOJO() map[string]interface{} {
	labels := []interface{}{}
	data := []interface{}{}
	dataset := map[string]interface{}{}
	datasets := []interface{}{}
	colors := []interface{}{}

	for _, c := range chartColors {
		colors = append(colors, c)
	}

	opByDate := map[string][]float32{}

	// Find op by Fund and add (+) all Operations
	for _, ops := range w.Ops {
		for _, op := range ops {
			date := op.GetDate().Format("02/01/2006")
			if _, ok := opByDate[date]; !ok {
				opByDate[date] = []float32{}
				opByDate[date] = append(opByDate[date], float32(0))
				opByDate[date] = append(opByDate[date], float32(0))
			}
			opByDate[date][0] += op.GetAmount()
			opByDate[date][1] += op.GetNav() * op.GetShares()
		}
	}

	labelsString := []string{"Input", "Worth"}
	for idx := 0; idx < 2; idx++ {
		for date, dayData := range opByDate {
			if idx == 0 {
				labels = append(labels, date)
			}
			data = append(data, dayData[idx])
		}

		dataset["label"] = labelsString[idx]
		dataset["backgroundColor"] = colors[idx*1]
		dataset["data"] = data

		datasets = append(datasets, dataset)

		data = []interface{}{}
		dataset = map[string]interface{}{}
	}

	return map[string]interface{}{
		"data": map[string]interface{}{
			"labels":   labels,
			"datasets": datasets,
		},
	}
}

// NewOperation returns a new operation
func NewOperation(date time.Time, fund string, amount, nav, shares float32) (Operation, error) {
	if !isValidPositive(amount) && !isValidPositive(nav) && !isValidPositive(shares) {
		return Operation{}, errors.New("operation values not valid")
	}

	return Operation{
		ID:       utils.NewULID(),
		FundInsi: fund,
		Date:     date,
		Amount:   amount,
		NAV:      nav,
		Shares:   shares,
	}, nil
}

// GetCopy creates a copy of the wallet
func (op Operation) GetCopy() Operation {
	return Operation{
		ID:       op.ID,
		FundInsi: op.FundInsi,
		Date:     op.Date,
		Amount:   op.Amount,
		NAV:      op.NAV,
		Shares:   op.Shares,
	}
}

// GetID returns the operation ID
func (op Operation) GetID() string {
	return op.ID
}

// GetInsi returns the INSI associated to the Fund where the operation was made
func (op Operation) GetInsi() string {
	return op.FundInsi
}

func (op *Operation) SetInsi(insi string) {
	// TODO: Validate insi
	op.FundInsi = insi
}

// GetDate returns the date where the operation as made
func (op Operation) GetDate() time.Time {
	return op.Date
}

// SetDate sets the operation's date
func (op *Operation) SetDate(t time.Time) {
	// TODO: Validate t
	op.Date = t
}

// GetAmount returns the input amount for the operation was carried out
func (op Operation) GetAmount() float32 {
	return op.Amount
}

// SetAmount sets the input for this operation
func (op *Operation) SetAmount(amount float32) {
	op.Amount = amount
}

// GetNav returns the Net Asset Value of the fund to which the operation
// was carried out
func (op Operation) GetNav() float32 {
	return op.NAV
}

// SetNav sets the Net Asset Value
func (op *Operation) SetNav(nav float32) {
	op.NAV = nav
}

// GetShares returns the number of shares for the operation was carried out
func (op Operation) GetShares() float32 {
	return op.Shares
}

// SetShares sets the shares field
func (op *Operation) SetShares(shares float32) {
	op.Shares = shares
}

func isValidName(name string) bool {
	filtered := strings.TrimSpace(name)
	return filtered != ""
}

func isValidPositive(value float32) bool {
	return value >= 0
}
