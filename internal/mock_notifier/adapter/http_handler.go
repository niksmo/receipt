package adapter

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/niksmo/receipt/internal/mock_notifier/core/domain"
	"github.com/niksmo/receipt/internal/mock_notifier/core/port"
	"github.com/niksmo/receipt/pkg/logger"
)

type SendMailHandler struct {
	log     logger.Logger
	service port.MessagePrinter
}

func RegisterMailReceiptHandler(
	log logger.Logger, mux *http.ServeMux, service port.MessagePrinter,
) {
	h := SendMailHandler{log, service}
	mux.HandleFunc("POST /v1/email", h.SendMail)
}

func (h SendMailHandler) SendMail(
	w http.ResponseWriter, r *http.Request,
) {
	const op = "SendMailHandler.SendMail"
	log := h.log.WithOp(op)

	var data SendEmail
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		errStr := "invalid JSON"
		http.Error(w, errStr, http.StatusBadRequest)
		log.Info().Err(err).Msg(errStr)
		return
	}

	msg, err := h.toDomain(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Info().Err(err).Msg(err.Error())
		return
	}

	msgID, err := h.service.PrintMessage(r.Context(), msg)
	if err != nil {
		http.Error(w, "", http.StatusServiceUnavailable)
		log.Error().Err(fmt.Errorf("%s: %w", op, err)).Msg("unexpected error")
		return
	}

	w.WriteHeader(http.StatusCreated)
	res := MessageCreated{MessageID: msgID.String()}
	_ = json.NewEncoder(w).Encode(res)
}

func (h SendMailHandler) toDomain(data SendEmail) (domain.Message, error) {
	if len(data.To) == 0 {
		return domain.Message{}, errors.New("empty field: 'to'")
	}

	msg := domain.Message{
		FromEmail: data.Sender.Email,
		ToEmail:   data.To[0].Email,
		Subject:   data.Subject,
		Content:   data.TextContent,
	}

	return msg, nil
}
