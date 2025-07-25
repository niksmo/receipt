package schema

import "time"

type Product struct {
	Name       string `json:"name"`
	Quantity   int    `json:"quantity"`
	UnitPrice  int    `json:"unit_price"`  // копейки
	TotalPrice int    `json:"total_price"` // копейки
	TaxRate    string `json:"tax_rate"`
	TaxValue   int    `json:"tax_value"` // копейки
}

type Receipt struct {
	Number             int       `json:"number"` // номер чека
	Date               time.Time `json:"date"`
	Organization       string    `json:"organization"`     // название организации
	PaymentAddress     string    `json:"payment_address"`  // адрес расчетов
	TaxpayerNumber     string    `json:"taxpayer_number"`  // ИНН
	TaxationType       string    `json:"taxation_type"`    // вид налогообложения
	CalculationSign    string    `json:"calculation_sign"` // признак расчета
	Products           []Product `json:"products"`
	CustomerEmail      string    `json:"customer_email"`
	FiscalDeviceNumber string    `json:"fiscal_device_number"` // ФН
	CashRegisterNumber string    `json:"cash_register_number"` // РН ККТ
	FiscalDocument     string    `json:"fiscal_document"`      // ФД
	FiscalAttribute    string    `json:"fiscal_attribute"`     // ФПД
}
