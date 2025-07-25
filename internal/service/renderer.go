package service

import (
	"bytes"
	_ "embed"
	"html/template"
	"strings"
	"time"

	"github.com/niksmo/receipt/internal/scheme"
)

const dateLayout = "01.02.06 15:04"

type tax struct {
	Rate  string
	Value float64
}

//go:embed receipt.template
var receiptTemplate string

func renderReciept(receipt scheme.Receipt) []byte {
	funcMap := template.FuncMap{
		"formatDate": formatDate,
		"upper":      strings.ToUpper,
		"lower":      strings.ToLower,
		"cost":       cost,
		"totalPrice": totalPrice,
		"taxes":      taxes,
	}

	t := template.Must(
		template.New("receipt").Funcs(funcMap).Parse(receiptTemplate))

	var b bytes.Buffer
	t.Execute(&b, receipt)
	return b.Bytes()
}

func formatDate(t time.Time) string {
	return t.Format(dateLayout)
}

func cost(price int) float64 {
	return float64(price) / 100
}

func totalPrice(products []scheme.Product) float64 {
	var sum int
	for _, p := range products {
		sum += p.TotalPrice
	}
	return float64(sum) / 100
}

func taxes(products []scheme.Product) []tax {
	var taxes []tax
	for _, p := range products {
		if p.TaxRate == "" {
			continue
		}
		t := tax{
			Rate:  p.TaxRate,
			Value: float64(p.TaxValue) / 100,
		}
		taxes = append(taxes, t)
	}
	return taxes
}
