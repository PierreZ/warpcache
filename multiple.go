package warpcache

// MultipleCache is watching other multiples GTS
type MultipleCache struct {
	cache
	pivot  string
	Values map[string]float64
}
