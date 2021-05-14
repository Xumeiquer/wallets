package models

import "reflect"

var (
	chartColors = []string{
		"#25CCF7", "#FD7272", "#54a0ff", "#00d2d3",
		"#1abc9c", "#2ecc71", "#3498db", "#9b59b6", "#34495e",
		"#16a085", "#27ae60", "#2980b9", "#8e44ad", "#2c3e50",
		"#f1c40f", "#e67e22", "#e74c3c", "#ecf0f1", "#95a5a6",
		"#f39c12", "#d35400", "#c0392b", "#bdc3c7", "#7f8c8d",
		"#55efc4", "#81ecec", "#74b9ff", "#a29bfe", "#dfe6e9",
		"#00b894", "#00cec9", "#0984e3", "#6c5ce7", "#ffeaa7",
		"#fab1a0", "#ff7675", "#fd79a8", "#fdcb6e", "#e17055",
		"#d63031", "#feca57", "#5f27cd", "#54a0ff", "#01a3a4",
	}
)

// CopyFundsMap creates a deep copy for Funds type
func CopyFundsMap(m Funds) Funds {
	res := make(Funds)
	for k, v := range m {
		tmp := FundData{}

		stSRC := reflect.ValueOf(&v).Elem()
		stDST := reflect.ValueOf(&tmp).Elem()

		for i := 0; i < stDST.NumField(); i++ {
			fDST := stDST.Field(i)
			if fDST.CanSet() {
				fDST.Set(stSRC.Field(i))
			}
		}
		res[k] = tmp

	}
	return res
}

// CopyOperationsMap creates a deep copy for Operations type
func CopyOperationsMap(ops Operations) Operations {
	res := make(Operations)

	for k, vs := range ops {
		res[k] = []Operation{}
		for _, v := range vs {
			tmp := Operation{}

			stSRC := reflect.ValueOf(&v).Elem()
			stDST := reflect.ValueOf(&tmp).Elem()

			for i := 0; i < stDST.NumField(); i++ {
				fDST := stDST.Field(i)
				if fDST.CanSet() {
					fDST.Set(stSRC.Field(i))
				}
			}
			res[k] = append(res[k], tmp)
		}
	}
	return res
}
