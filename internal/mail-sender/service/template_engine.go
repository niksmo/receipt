package service

import (
	"bytes"
	_ "embed"
	"html/template"
	"strings"

	"github.com/niksmo/receipt/internal/mail-sender/domain"
)

//go:embed text.template
var receiptTemplate string

type TextTemplateEngine struct {
	template *template.Template
}

func NewTextTemplateEngine() (TextTemplateEngine, error) {
	t, err := template.New("receipt").Funcs(
		template.FuncMap{
			"upper": strings.ToUpper,
			"lower": strings.ToLower,
		},
	).Parse(receiptTemplate)
	if err != nil {
		return TextTemplateEngine{}, err
	}
	return TextTemplateEngine{t}, nil
}

func (t TextTemplateEngine) TemplateToString(r domain.Receipt) string {
	var b bytes.Buffer
	t.template.Execute(&b, r)
	return b.String()
}
