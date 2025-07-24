package schema

import "time"

// цена и налог в копейках, центах и т.п.
type Product struct {
	Name         string
	Quantity     int
	PricePerUnit int
	TotalPrice   int
	TaxRate      string
	TaxValue     int
}

type Receipt struct {
	Number             int // номер чека
	Date               time.Time
	Organization       string // название организации
	PaymentAddress     string // адрес расчетов
	TaxpayerNumber     string // ИНН
	TaxationType       string // вид налогообложения
	CalculationSign    string // признак расчета
	Products           []Product
	BuyerEmail         string
	FiscalDeviceNumber string // ФН
	CashRegisterNumber string // РН ККТ
	FiscalDocument     string // ФД
	FiscalAttribute    string // ФПД
}
