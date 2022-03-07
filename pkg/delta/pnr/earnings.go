package pnr

import (
	"fmt"
	"strconv"
)

func asFloat(n string) float64 {
	s, _ := strconv.ParseFloat(n, 64)
	return s
}

func asString(n float64) string {
	return fmt.Sprintf("%.2f", n)
}

func estimateMQD(pnr *PNR) string {
	// We can't calculate MQD very easily for non-USD.
	if pnr.Fare.BaseCurrencyCode != "USD" || pnr.Fare.TotalCurrencyCode != "USD" {
		return "unknown currency"
	}

	total := float64(0)
	total += asFloat(pnr.Fare.BaseFare)

	for _, row := range pnr.Fare.TaxRows {
		if row.CarrierImposedFee && row.Currency == "USD" {
			total += asFloat(row.Amount)
		} else if row.Currency != "USD" {
			return "unknown currency"
		}
	}

	return fmt.Sprintf("%s %s", asString(total), pnr.Fare.BaseCurrencyCode)
}
