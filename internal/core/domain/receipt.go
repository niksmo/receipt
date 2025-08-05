package domain

import (
	"cmp"
	"slices"
	"time"

	"github.com/google/uuid"
)

type Tax struct {
	TaxRate  string
	TaxValue int
}

type Product struct {
	Name       string
	Quantity   int
	UnitPrice  int
	TotalPrice int
	TaxRate    string
	TaxValue   int
}

type Receipt struct {
	UUID               string
	Number             int
	Date               time.Time
	Organization       string
	PaymentAddress     string
	TaxpayerNumber     string
	TaxationType       string
	CalculationSign    string
	CustomerEmail      string
	FiscalDeviceNumber string
	CashRegisterNumber string
	FiscalDocument     string
	FiscalAttribute    string
	Products           []Product
}

func NewReceipt() Receipt {
	return Receipt{UUID: uuid.NewString()}
}

func (r *Receipt) TotalPrice() int {
	var total int
	for _, p := range r.Products {
		total += p.TotalPrice
	}
	return total
}

func (r *Receipt) TotalTax() []Tax {
	index := make(map[string]int)

	for _, p := range r.Products {
		if p.TaxRate != "" {
			index[p.TaxRate] += p.TaxValue
		}
	}

	tax := make([]Tax, 0, len(index))
	for r, v := range index {
		tax = append(tax, Tax{TaxRate: r, TaxValue: v})
	}

	slices.SortFunc(tax, func(a, b Tax) int {
		return cmp.Compare(a.TaxRate, b.TaxRate)
	})

	return tax
}
