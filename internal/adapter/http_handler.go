package adapter

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/niksmo/receipt/internal/core/domain"
	"github.com/niksmo/receipt/internal/core/port"
	"github.com/niksmo/receipt/pkg/logger"
)

type MailReceiptHandler struct {
	log     logger.Logger
	service port.EventSaver
}

func RegisterMailReceiptHandler(
	log logger.Logger, mux *http.ServeMux, service port.EventSaver,
) {
	h := MailReceiptHandler{log, service}
	mux.HandleFunc("POST /v1/receipt", h.SendReceiptToMail)
}

func (h MailReceiptHandler) SendReceiptToMail(
	w http.ResponseWriter, r *http.Request,
) {
	const op = "MailReceiptHandler.SendReceiptToMail"
	log := h.log.WithOp(op)

	if r.Header.Get("Content-Type") != "application/json" {
		errStr := "invalid media type"
		http.Error(w, errStr, http.StatusUnsupportedMediaType)
		log.Info().Msg(errStr)
		return
	}

	var data Receipt
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		errStr := "invalid json"
		http.Error(w, errStr, http.StatusBadRequest)
		log.Info().Err(err).Msg(errStr)
		return
	}

	receipt := httpReceiptToDomain(data)
	err = h.service.SaveEvent(r.Context(), receipt)
	if err != nil {
		http.Error(w, "", http.StatusServiceUnavailable)
		log.Error().Err(fmt.Errorf("%s: %w", op, err)).Msg("unexpected error")
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Accept"))
}

func httpReceiptToDomain(data Receipt) domain.Receipt {
	r := domain.NewReceipt()
	r.Number = data.Number
	r.Date = data.Date
	r.Organization = data.Organization
	r.PaymentAddress = data.PaymentAddress
	r.TaxpayerNumber = data.TaxpayerNumber
	r.TaxationType = data.TaxationType
	r.CalculationSign = data.CalculationSign
	r.CustomerEmail = data.CustomerEmail
	r.FiscalDeviceNumber = data.FiscalDeviceNumber
	r.CashRegisterNumber = data.CashRegisterNumber
	r.FiscalDocument = data.FiscalDocument
	r.FiscalAttribute = data.FiscalAttribute

	for _, pd := range data.Products {
		r.Products = append(
			r.Products,
			domain.Product{
				Name:       pd.Name,
				Quantity:   pd.Quantity,
				UnitPrice:  pd.UnitPrice,
				TotalPrice: pd.TotalPrice,
				TaxRate:    pd.TaxRate,
				TaxValue:   pd.TaxValue,
			},
		)
	}
	return r
}
