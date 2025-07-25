package httpapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/niksmo/receipt/internal/scheme"
	"github.com/niksmo/receipt/pkg/logger"
)

type ReceiptSender interface {
	SendReceipt(context.Context, scheme.Receipt) error
}

type MailReceiptHandler struct {
	log     logger.Logger
	service ReceiptSender
}

func RegisterMailReceiptHandler(
	log logger.Logger, mux *http.ServeMux, service ReceiptSender,
) {
	h := MailReceiptHandler{log, service}
	mux.HandleFunc("POST /v1/receipt", h.receiptPOSTv1)
}

func (h MailReceiptHandler) receiptPOSTv1(
	w http.ResponseWriter, r *http.Request,
) {
	const op = "MailReceiptHandler.receiptPOSTv1"
	log := h.log.WithOp(op)

	var receipt scheme.Receipt

	if r.Header.Get("Content-Type") != "application/json" {
		errStr := "invalid media type"
		http.Error(w, errStr, http.StatusUnsupportedMediaType)
		log.Info().Msg(errStr)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&receipt)
	if err != nil {
		errStr := "invalid json data"
		http.Error(w, "invalid json data", http.StatusBadRequest)
		log.Info().Err(err).Msg(errStr)
		return
	}

	err = h.service.SendReceipt(r.Context(), receipt)
	if err != nil {
		http.Error(w, "", http.StatusServiceUnavailable)
		log.Error().Err(fmt.Errorf("%s: %w", op, err)).Msg("unexpected error")
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Info().Str("customer", receipt.CustomerEmail).Msg("sent")
}
