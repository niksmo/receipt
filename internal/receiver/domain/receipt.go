package domain

import (
	"time"

	"github.com/google/uuid"
)

const dateLayout = "01.02.06 15:04"

type ProductData struct {
	Name       string
	Quantity   int
	UnitPrice  int
	TotalPrice int
	TaxRate    string
	TaxValue   int
}

type ReceiptData struct {
	Number             int
	Date               time.Time
	Organization       string
	PaymentAddress     string
	TaxpayerNumber     string
	TaxationType       string
	CalculationSign    string
	Products           []ProductData
	CustomerEmail      string
	FiscalDeviceNumber string
	CashRegisterNumber string
	FiscalDocument     string
	FiscalAttribute    string
}

type Tax struct {
	TaxRate  string
	TaxValue float64
}

type Product struct {
	Name       string
	Quantity   int
	UnitPrice  float64
	TotalPrice float64
	Tax
}

type Receipt struct {
	UUID               string
	Number             int
	Date               string
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
	TotalPrice         float64
	TotalTax           []Tax
	Products           []Product
}

func NewReceipt(rd ReceiptData) (Receipt, error) {
	r := Receipt{
		UUID:               uuid.NewString(),
		Organization:       rd.Organization,
		PaymentAddress:     rd.PaymentAddress,
		TaxpayerNumber:     rd.TaxpayerNumber,
		TaxationType:       rd.TaxationType,
		CalculationSign:    rd.CalculationSign,
		CustomerEmail:      rd.CustomerEmail,
		FiscalDeviceNumber: rd.FiscalDeviceNumber,
		FiscalAttribute:    rd.FiscalAttribute,
	}
	r.setDate(rd.Date)
	r.setTotalPrice(rd.Products)
	r.setProducts(rd.Products)
	r.setTaxes(rd.Products)
	return r, nil
}

func (r *Receipt) setDate(t time.Time) {
	r.Date = t.Format(dateLayout)
}

func (r *Receipt) setTotalPrice(productsData []ProductData) {
	var sum int
	for _, p := range productsData {
		sum += p.TotalPrice
	}
	r.TotalPrice = cost(sum)
}

func (r *Receipt) setProducts(productsData []ProductData) {
	var products []Product
	for _, p := range productsData {
		products = append(products, Product{
			Name:       p.Name,
			Quantity:   p.Quantity,
			UnitPrice:  cost(p.UnitPrice),
			TotalPrice: cost(p.TotalPrice),
			Tax: Tax{
				TaxRate:  p.TaxRate,
				TaxValue: cost(p.TaxValue),
			},
		})
	}
	r.Products = products
}

func (r *Receipt) setTaxes(productsData []ProductData) {
	taxesTmp := make(map[string]int)
	for _, p := range productsData {
		tr := p.TaxRate
		taxesTmp[tr] += p.TaxValue
	}

	taxes := make([]Tax, len(taxesTmp))
	for rate, v := range taxesTmp {
		taxes = append(taxes, Tax{TaxRate: rate, TaxValue: cost(v)})
	}
	r.TotalTax = taxes
}

func cost(price int) float64 {
	return float64(price) / 100
}
