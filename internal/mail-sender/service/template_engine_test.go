//go:build !integration

package service_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/niksmo/receipt/internal/mail-sender/domain"
	"github.com/niksmo/receipt/internal/mail-sender/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRenderer(t *testing.T) {
	receipt := domain.Receipt{
		UUID:               uuid.NewString(),
		Number:             1234,
		Date:               "07.25.25 14:40",
		Organization:       "ООО Ромашка",
		PaymentAddress:     "г. Москва, ул. Правды, д. 1",
		TaxpayerNumber:     "7745123451234",
		TaxationType:       "ОСН",
		CalculationSign:    "приход",
		CustomerEmail:      "Happy_Customer@mail.ru",
		FiscalDeviceNumber: "7380440801479592",
		CashRegisterNumber: "0007768750034436",
		FiscalDocument:     "16415",
		FiscalAttribute:    "1805600812",
		Products: []domain.Product{
			{
				Name:       "тапки синие размер 42",
				Quantity:   1,
				UnitPrice:  230,
				TotalPrice: 230,
			},
			{
				Name:       "очки солнцезащитные",
				Quantity:   1,
				UnitPrice:  6345,
				TotalPrice: 6345,
				Tax: domain.Tax{
					TaxRate:  "5/105",
					TaxValue: 317.25,
				},
			},
			{
				Name:       "мыло душистое",
				Quantity:   5,
				UnitPrice:  80,
				TotalPrice: 400,
				Tax: domain.Tax{
					TaxRate:  "20",
					TaxValue: 80,
				},
			},
			{
				Name:       "туалетная вода",
				Quantity:   1,
				UnitPrice:  6340,
				TotalPrice: 6340,
				Tax: domain.Tax{
					TaxRate:  "20",
					TaxValue: 1268,
				},
			},
		},
		TotalPrice: 13315.00,
		TotalTax: []domain.Tax{
			{TaxRate: "5/105", TaxValue: 317.25},
			{TaxRate: "20", TaxValue: 1348.00},
		},
	}

	expected := `Кассовый чек № 1234
07.25.25 14:40

ООО Ромашка
г. Москва, ул. Правды, д. 1
ИНН 7745123451234
Вид налогообложения: ОСН

ПРИХОД

тапки синие размер 42
1 x 230.00
=230.00
без НДС

очки солнцезащитные
1 x 6345.00
=6345.00
в т.ч. НДС 5/105
= 317.25

мыло душистое
5 x 80.00
=400.00
в т.ч. НДС 20
= 80.00

туалетная вода
1 x 6340.00
=6340.00
в т.ч. НДС 20
= 1268.00

--
ИТОГ =13315.00
в т.ч. НДС 5/105 =317.25
в т.ч. НДС 20 =1348.00
Безналичными =13315.00

Электронный адрес покупателя
happy_customer@mail.ru

ФН: 7380440801479592
РН ККТ: 0007768750034436
ФД: 16415
ФПД: 1805600812
`

	templateEngine, err := service.NewTextTemplateEngine()
	require.NoError(t, err)

	actual := templateEngine.TemplateToString(receipt)
	assert.Equal(t, expected, actual)
}
