package ports

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/niksmo/receipt/internal/receiver/domain"
	"github.com/niksmo/receipt/pkg/logger"
)

type ReceiptSender interface {
	SendReceipt(context.Context, domain.Receipt) error
}

type MailReceiptHandler struct {
	log     logger.Logger
	service ReceiptSender
}

func RegisterMailReceiptHandler(
	log logger.Logger, mux *http.ServeMux, service ReceiptSender,
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

	var receiptData Receipt
	err := json.NewDecoder(r.Body).Decode(&receiptData)
	if err != nil {
		errStr := "invalid json"
		http.Error(w, errStr, http.StatusBadRequest)
		log.Info().Err(err).Msg(errStr)
		return
	}

	receipt, err := httpReceiptToDomain(receiptData)
	if err != nil {
		errStr := "invalid data"
		http.Error(w, errStr, http.StatusBadRequest)
		log.Info().Err(err).Msg(errStr)
		return
	}

	err = h.service.SendReceipt(r.Context(), receipt)
	if err != nil {
		http.Error(w, "", http.StatusServiceUnavailable)
		log.Error().Err(fmt.Errorf("%s: %w", op, err)).Msg("unexpected error")
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Accept"))
}

func httpReceiptToDomain(reqData Receipt) (domain.Receipt, error) {
	var products []domain.ProductData
	for _, p := range reqData.Products {
		products = append(products, domain.ProductData{
			Name:       p.Name,
			Quantity:   p.Quantity,
			UnitPrice:  p.UnitPrice,
			TotalPrice: p.TotalPrice,
			TaxRate:    p.TaxRate,
			TaxValue:   p.TaxValue,
		})
	}
	return domain.NewReceipt(domain.ReceiptData{
		Number:             reqData.Number,
		Date:               reqData.Date,
		Organization:       reqData.Organization,
		PaymentAddress:     reqData.Organization,
		TaxpayerNumber:     reqData.TaxpayerNumber,
		TaxationType:       reqData.TaxationType,
		CalculationSign:    reqData.CalculationSign,
		CustomerEmail:      reqData.CustomerEmail,
		FiscalDeviceNumber: reqData.FiscalDeviceNumber,
		FiscalDocument:     reqData.FiscalDocument,
		CashRegisterNumber: reqData.CashRegisterNumber,
		FiscalAttribute:    reqData.FiscalAttribute,
		Products:           products,
	})
}
