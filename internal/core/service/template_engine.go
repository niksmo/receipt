package service

import (
	"bytes"
	_ "embed"
	"html/template"
	"strings"
	"time"

	"github.com/niksmo/receipt/internal/core/domain"
)

const DateLayout = "01.02.06 15:04"

//go:embed text.template
var receiptTemplate string

type TextTemplateEngine struct {
	template *template.Template
}

func NewReceiptTemplateEngine() TextTemplateEngine {
	funcMap := template.FuncMap{
		"formatDate": formatDate,
		"upper":      strings.ToUpper,
		"lower":      strings.ToLower,
		"cost":       cost,
	}

	t, err := template.New("receipt").Funcs(funcMap).Parse(receiptTemplate)
	if err != nil {
		panic(err)
	}

	return TextTemplateEngine{t}
}

func (t TextTemplateEngine) ToText(r *domain.Receipt) string {
	var b bytes.Buffer
	t.template.Execute(&b, r)
	return b.String()
}

func formatDate(d time.Time) string {
	return d.Format(DateLayout)
}

func cost(price int) float64 {
	return float64(price) / 100
}
