package service

import (
	"bytes"
	_ "embed"
	"html/template"
	"strings"
	"time"

	"github.com/niksmo/receipt/internal/schema"
)

const dateLayout = "01.02.06 15:04"

type tax struct {
	Rate  string
	Value float64
}

type renderer struct {
	Taxes []tax
	schema.Receipt
}

func (r *renderer) ProductTax(p schema.Product) float64 {
	t := tax{
		Rate:  p.TaxRate,
		Value: float64(p.TaxValue) / 100,
	}
	r.Taxes = append(r.Taxes, t)
	return t.Value
}

func (r *renderer) TotalPrice() float64 {
	var sum int
	for _, p := range r.Receipt.Products {
		sum += p.TotalPrice
	}
	return float64(sum) / 100
}

//go:embed receipt.template
var receiptTemplate string

func renderReciept(receipt schema.Receipt) []byte {
	funcMap := template.FuncMap{
		"formatDate": formatDate,
		"upper":      strings.ToUpper,
		"lower":      strings.ToLower,
		"cost":       cost,
	}

	t := template.Must(
		template.New("receipt").Funcs(funcMap).Parse(receiptTemplate))

	var b bytes.Buffer
	renderer := renderer{Receipt: receipt}
	t.Execute(&b, &renderer)
	return b.Bytes()
}

func formatDate(t time.Time) string {
	return t.Format(dateLayout)
}

func cost(price int) float64 {
	return float64(price) / 100
}
