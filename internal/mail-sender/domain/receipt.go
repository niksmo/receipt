package domain

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
