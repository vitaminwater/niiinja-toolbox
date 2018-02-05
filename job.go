package niiinja

import (
	"errors"
	"reflect"

	tachikoma "github.com/vitaminwater/tachikoma-toolbox"
	"github.com/vitaminwater/tachikoma-toolbox/timeseries"
)

// Uses the fixer API to convert currencies

type CurrencyConvertJob struct {
	timeseries.Job
	from string
	to   string
}

func (j CurrencyConvertJob) Run(i interface{}) error {
	if t, ok := i.(CurrencyTicker); ok == false {
		return errors.New("CurrencyConvertJob required a CurrencyTicker object")
	} else if t.GetBase() != j.from && t.GetCounter() != j.from {
		return nil
	}
	return j.Job.Run(i)
}

func ConversionSelector(from, to string, s tachikoma.Selector) tachikoma.Selector {
	v := GetSymbolValue(from, to)
	return func(o interface{}) interface{} {
		var ct CurrencyTicker
		var ok bool
		if ct, ok = o.(CurrencyTicker); ok == false {
			tachikoma.Fatal(errors.New("CurrencyConvertJob required a CurrencyTicker object"))
		}
		o = s(o)
		if reflect.TypeOf(o).Kind() != reflect.Float64 {
			tachikoma.Fatal(errors.New("ConversionSelector requires a float64 value"))
		}
		f := reflect.ValueOf(o).Float()
		var mult float64 = 1
		if ct.GetBase() == from {
			mult = 1 / v
		} else if ct.GetCounter() == from {
			mult = v
		}
		return f * mult
	}
}

func NewCurrencyConvertJob(name, from, to string, timeserie timeseries.Timeserie, labels timeseries.Labels, selector tachikoma.Selector) tachikoma.Job {
	labels = labels.Clone()
	labels["base"] = labelTickerFn(func(t CurrencyTicker) string {
		if t.GetBase() == from {
			return to
		}
		return t.GetBase()
	})
	labels["counter"] = labelTickerFn(func(t CurrencyTicker) string {
		if t.GetCounter() == from {
			return to
		}
		return t.GetCounter()
	})
	timeserie.Name = name
	j := timeseries.NewJob(timeserie, labels, ConversionSelector(from, to, selector))

	cc := CurrencyConvertJob{
		Job:  j,
		from: from,
		to:   to,
	}
	return cc
}
