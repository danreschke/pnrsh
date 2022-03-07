package pnr

import "strings"

func convertRemarks(res RetrievePnrResponse, pnr *PNR) {
	for _, remark := range res.TripsResponse.Journey.Pnr.Remarks.DomainObjectList.DomainObject {
		pnr.Remarks = append(pnr.Remarks, Remark{
			FreeFormText: remark.FreeFormText,
			RemarkType:   remark.RemarkType,
		})
	}
}

func convertFlights(res RetrievePnrResponse, pnr *PNR) {
	for _, itin := range res.TripsResponse.Journey.Pnr.Itineraries.DomainObjectList.DomainObject {
		for _, flight := range itin.Flights.DomainObjectList.DomainObject {
			pnr.Flights = append(pnr.Flights, Flight{
				OriginAirportCode:      flight.Origin.Code,
				DestinationAirportCode: flight.Destination.Code,
				CurrentActionCode:      flight.CurrentActionCode,
				PreviousActionCode:     flight.PreviousActionCode,
				Distance:               flight.Distance,
				Status:                 flight.Status,
				MarketingAirlineCode:   flight.MarketingAirlineCode,
				OperatingAirlineCode:   flight.OperatingAirlineCode,
				UpgradeStatus:          flight.UpgradeStatus,
			})
		}
	}
}

func convertPassengers(res RetrievePnrResponse, pnr *PNR) {
	for _, pax := range res.TripsResponse.Journey.Pnr.Passengers.DomainObjectList.DomainObject {
		convertedPax := Passenger{
			Name:       pax.Name.FirstName + " " + pax.Name.LastName,
			CustomerID: pax.CustomerId,
			CheckedIn:  pax.CheckedIn != "false",
		}

		for _, ssr := range pax.Ssrs.DomainObjectList.DomainObject {
			if ssr.Code == "FQTU" && (strings.Contains(ssr.Remarks.Remark, "OU") || strings.Contains(ssr.Remarks.Remark, "SU")) {
				continue
			}

			convertedPax.SSRs = append(convertedPax.SSRs, SSR{
				Remark: ssr.Remarks.Remark,
				Type:   ssr.Code,
			})
		}
		pnr.Passengers = append(pnr.Passengers, convertedPax)
	}
}

func convertFlags(res RetrievePnrResponse, pnr *PNR) {
	for _, flag := range res.TripsResponse.Journey.Pnr.PnrFlags.DomainObjectList.DomainObject {
		if flag.Name == "" || flag.Value == "" {
			continue
		}

		pnr.Flags = append(pnr.Flags, Flag{
			Name:  flag.Name,
			Value: flag.Value,
		})
	}
}

func convertTickets(res RetrievePnrResponse, pnr *PNR) {
	for _, pax := range res.TripsResponse.Journey.Pnr.Passengers.DomainObjectList.DomainObject {
		for _, ticket := range pax.Tickets.DomainObjectList.DomainObject {
			pnr.Tickets = append(pnr.Tickets, Ticket{
				Number:                 ticket.Number,
				ExpirationDate:         ticket.ExpirationDate,
				IssueDate:              ticket.IssueDate,
				Status:                 ticket.Status,
				PassengerName:          pax.Name.FirstName + " " + pax.Name.LastName,
				NumCoupons:             uint64(len(ticket.TicketCoupons.DomainObjectList.DomainObject)),
				ValidatedAgainstCoupon: couponsMatchFlights(res, ticket.Number),
			})
		}
	}
}

func convertFare(res RetrievePnrResponse, pnr *PNR) {
	fare := res.TripsResponse.Journey.Pnr.TotalFare

	convertedFare := Fare{
		BaseCurrencyCode:  fare.BaseCurrencyCode,
		BaseFare:          fare.BaseFare,
		TotalTax:          fare.TotalTax,
		TotalCurrencyCode: fare.TotalCurrencyCode,
		TotalFare:         fare.TotalFare,
	}

	for _, row := range fare.TaxBreakDownList.FareFaxTable {
		convertedFare.TaxRows = append(convertedFare.TaxRows, struct {
			TaxType           string
			Amount            string
			Currency          string
			CarrierImposedFee bool
		}{
			row.TaxType,
			row.Amount,
			row.Currency,
			row.CarrierImposedFee != "false",
		})
	}

	pnr.Fare = convertedFare
	pnr.Fare.EstimatedMQD = estimateMQD(pnr)
}
