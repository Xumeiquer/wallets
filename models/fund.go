package models

import "time"

// FundData represents the dynamic data for the Fund
type FundData struct {
	INSI        string     `json:"insi,omitempty"`
	Name        string     `json:"name,omitempty"`
	Nav         float32    `json:"nav,omitempty"`
	Date        time.Time  `json:"date,omitempty"`
	DayChange   [2]float32 `json:"day_change,omitempty"`
	TimeDelayed string     `json:"time_delayed,omitempty"`
	OneMonthS   string     `json:"one_month_s,omitempty"`
	ThreeMonthS string     `json:"three_month_s,omitempty"`
	SixMonthS   string     `json:"six_month_s,omitempty"`
	OneYearS    string     `json:"one_year_s,omitempty"`

	OneDay     float32 `json:"one_day,omitempty"`
	OneWeek    float32 `json:"one_week,omitempty"`
	OneMonth   float32 `json:"one_month,omitempty"`
	ThreeMonth float32 `json:"three_month,omitempty"`
	OneYear    float32 `json:"one_year,omitempty"`
	ThreeYears float32 `json:"three_years,omitempty"`
	FiveYears  float32 `json:"five_years,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
}

// NewFundData returns a new empty FundData
func NewFundData(insi string) FundData {
	return FundData{
		INSI:      insi,
		CreatedAt: time.Now(),
	}
}

// GetINSI returns the INSI number for the Fund
func (f FundData) GetINSI() string {
	return f.INSI
}

// GetName returns the Fund name
func (f FundData) GetName() string {
	return f.Name
}

// GetNav returns the net asset value
func (f FundData) GetNav() float32 {
	return f.Nav
}

// GetDate returns the date which the operation was carried out
func (f FundData) GetDate() time.Time {
	return f.Date
}

// Predictions returns the profit predictions
func (f FundData) Predictions() []float32 {
	res := []float32{}
	res = append(res, f.OneDay)
	res = append(res, f.OneWeek)
	res = append(res, f.OneMonth)
	res = append(res, f.ThreeMonth)
	res = append(res, f.OneYear)
	res = append(res, f.ThreeYears)
	res = append(res, f.FiveYears)
	return res
}
