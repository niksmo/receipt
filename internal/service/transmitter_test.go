package service

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSMTPTransmitter(t *testing.T) {
	login := os.Getenv("TEST_RECEIPT_LOGIN")
	password := os.Getenv("TEST_RECEIPT_PASSWORD")
	host := os.Getenv("TEST_RECEIPT_HOST")
	port := os.Getenv("TEST_RECEIPT_PORT")
	to := os.Getenv("TEST_RECEIPT_TO")

	for _, v := range []string{login, password, host, port, to} {
		if v == "" {
			t.Skip("ENV variables not set")
		}
	}

	smtpTransmitter := NewSMTPTransmitter(login, password, host, port)
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	data := []byte(`Кассовый чек № 1234
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

--
ИТОГ =6975.00
в т.ч. НДС 5/105 =317.25
в т.ч. НДС 20 =80.00
Безналичными =6975.00

Электронный адрес покупателя
happy_customer@mail.ru

ФН: 7380440801479592
РН ККТ: 0007768750034436
ФД: 16415
ФПД: 1805600812
`)
	subject := "test receipt"
	err := smtpTransmitter.Send(ctx, to, subject, data)
	assert.NoError(t, err)
}
