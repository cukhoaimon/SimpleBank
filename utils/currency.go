package utils

const (
	USD = "USD"
	CAD = "CAD"
	EUR = "EUR"
	VND = "VND"
)

// IsSupportedCurrency return true if the currency is support
// otherwise is false
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, CAD, EUR, VND:
		return true
	}
	return false
}
