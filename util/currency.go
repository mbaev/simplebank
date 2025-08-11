package util

const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
)

var SupportedCurrencies = []string{EUR, USD, CAD}

func IsSupportedCurrency(currency string) bool {
	for _, supportedCurrency := range SupportedCurrencies {
		if currency == supportedCurrency {
			return true
		}
	}
	return false
}
