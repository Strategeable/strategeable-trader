package types

import "fmt"

type Symbol struct {
	BaseAsset          string
	BaseAssetPrecision int

	QuoteAsset          string
	QuoteAssetPrecision int

	// TODO: Add filters (LOT_SIZE, MIN_NOTIONAL, etc)
}

func NewSymbol(baseAsset string, baseAssetPrecision int, quoteAsset string, quoteAssetPrecision int) Symbol {
	return Symbol{
		BaseAsset:           baseAsset,
		BaseAssetPrecision:  baseAssetPrecision,
		QuoteAsset:          quoteAsset,
		QuoteAssetPrecision: quoteAssetPrecision,
	}
}

func (s *Symbol) String() string {
	return fmt.Sprintf("%s/%s", s.BaseAsset, s.QuoteAsset)
}

func (s *Symbol) MinQuoteSize() float64 {
	return 1
}
