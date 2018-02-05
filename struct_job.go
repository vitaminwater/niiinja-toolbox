package niiinja

import (
	"errors"
	"reflect"
	"strings"

	tachikoma "github.com/vitaminwater/tachikoma-toolbox"
	"github.com/vitaminwater/tachikoma-toolbox/timeseries"
)

type CurrencyTicker interface {
	GetBase() string
	GetCounter() string
}

type LabelTickerFn func(CurrencyTicker) string

func labelTickerFn(fn LabelTickerFn) timeseries.LabelFn {
	return func(f string, o interface{}) string {
		if t, ok := o.(CurrencyTicker); ok == true {
			return fn(t)
		}
		tachikoma.Fatal(errors.New("Ooops does not implement the CurrencyTicker inteface"))
		return ""
	}
}

func AddTickerLabels(labels timeseries.Labels) timeseries.Labels {
	labels["base"] = labelTickerFn(func(t CurrencyTicker) string {
		return t.GetBase()
	})
	labels["counter"] = labelTickerFn(func(t CurrencyTicker) string {
		return t.GetCounter()
	})
	return labels
}

func TickerJobGenerator(timeserie timeseries.Timeserie, labels timeseries.Labels, f reflect.StructField, selector tachikoma.Selector) []tachikoma.Job {
	jobs := make([]tachikoma.Job, 0)

	j := timeseries.NewJob(timeserie, labels, selector)
	jobs = append(jobs, j)

	if tag, ok := f.Tag.Lookup("convert"); ok == true {
		convs := strings.Split(tag, ",")
		for _, conv := range convs {
			cs := strings.Split(conv, ":")
			j := NewCurrencyConvertJob(cs[0], cs[1], timeserie, labels, selector)
			jobs = append(jobs, j)
		}
	}

	return jobs
}
