package niiinja

import (
	"fmt"

	tachikoma "github.com/vitaminwater/tachikoma-toolbox"
)

const USD_PRICE_URL = "https://api.fixer.io/latest?symbols=%s&base=%s"

type fixerAPI struct {
	Rates struct {
		USD float64
	}
}

func GetSymbolValue(from, to string) float64 {
	d := fixerAPI{}
	tachikoma.GetJSON(fmt.Sprintf(USD_PRICE_URL, to, from), &d)
	return d.Rates.USD
}
