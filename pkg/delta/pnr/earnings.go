package pnr

import (
	"fmt"
	"net/url"
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

func generateSmcalcRoute(pnr *PNR) (out string) {
	var fallbackClass string
	if len(pnr.Fare.FareBasisCode) > 0 {
		fallbackClass = string(pnr.Fare.FareBasisCode[0])
	}

	for idx, flight := range pnr.Flights {
		class := flight.ClassOfService

		if len(class) != 1 && len(fallbackClass) == 1 {
			class = fallbackClass
		} else if len(class) != 1 {
			class = "V"
		}

		if idx > 0 {
			out += fmt.Sprintf("-%s.%s-%s", flight.MarketingAirlineCode, class, flight.DestinationAirportCode)
		} else {
			out += fmt.Sprintf("%s-%s.%s-%s", flight.OriginAirportCode, flight.MarketingAirlineCode, class, flight.DestinationAirportCode)
		}
	}

	return out
}

func generateSmcalcLink(pnr *PNR) string {
	route := generateSmcalcRoute(pnr)
	return fmt.Sprintf("https://fly.qux.us/smcalc/dist.php?route=%s&start_mqm=0&start_rdm=0&default_fare=T&default_carrier=DL&elite=peon", url.QueryEscape(route))
}
