//go:build !integration

package service_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/niksmo/receipt/internal/receipt_service/core/domain"
	"github.com/niksmo/receipt/internal/receipt_service/core/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRenderer(t *testing.T) {
	date, err := time.Parse(service.DateLayout, "07.25.25 14:40")
	require.NoError(t, err)

	receipt := domain.Receipt{
		UUID:               uuid.NewString(),
		Number:             1234,
		Date:               date,
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
				Name:       "тапочки синие размер 42",
				Quantity:   1,
				UnitPrice:  23000,
				TotalPrice: 23000,
			},
			{
				Name:       "очки солнцезащитные",
				Quantity:   1,
				UnitPrice:  634500,
				TotalPrice: 634500,
				TaxRate:    "5/105",
				TaxValue:   31725,
			},
			{
				Name:       "мыло душистое",
				Quantity:   5,
				UnitPrice:  8000,
				TotalPrice: 40000,
				TaxRate:    "20",
				TaxValue:   8000,
			},
			{
				Name:       "туалетная вода",
				Quantity:   1,
				UnitPrice:  634000,
				TotalPrice: 634000,
				TaxRate:    "20",
				TaxValue:   126800,
			},
		},
	}

	expected := `Кассовый чек № 1234
07.25.25 14:40

ООО Ромашка
г. Москва, ул. Правды, д. 1
ИНН 7745123451234
Вид налогообложения: ОСН

ПРИХОД

тапочки синие размер 42
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
в т.ч. НДС 20 =1348.00
в т.ч. НДС 5/105 =317.25
Безналичными =13315.00

Электронный адрес покупателя
happy_customer@mail.ru

ФН: 7380440801479592
РН ККТ: 0007768750034436
ФД: 16415
ФПД: 1805600812
`

	templateEngine := service.NewReceiptTemplateEngine()
	actual := templateEngine.ToText(&receipt)
	assert.Equal(t, expected, actual)
}
