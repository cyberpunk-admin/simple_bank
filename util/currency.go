package util

// contain all support currencies
const (
	USD = "USD"
	EUR = "EUR"
	RMB = "RMB"
)

// IsSupport return if the currency is support
func IsSupport(currency string) bool {
	switch currency {
	case USD, EUR, RMB:
		return true
	}
	return false
}
