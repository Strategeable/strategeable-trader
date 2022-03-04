package types

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
