package adapter

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/niksmo/receipt/internal/mock_notifier/core/domain"
	"github.com/niksmo/receipt/internal/mock_notifier/core/port"
	"github.com/niksmo/receipt/pkg/logger"
	"golang.org/x/time/rate"
)

type SendMailHandler struct {
	log     logger.Logger
	service port.MessagePrinter
	limiter *rate.Limiter
}

func RegisterSendMailHandler(
	log logger.Logger,
	mux *http.ServeMux,
	service port.MessagePrinter,
	limit int,
) {
	limiter := createLimiter(limit)
	h := SendMailHandler{log, service, limiter}
	mux.HandleFunc("POST /v1/email", h.SendMail)
}

func (h SendMailHandler) SendMail(
	w http.ResponseWriter, r *http.Request,
) {
	const op = "SendMailHandler.SendMail"
	log := h.log.WithOp(op)

	if !h.allow(w) {
		log.Info().Msg("request limited")
		return
	}

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

func (h SendMailHandler) allow(w http.ResponseWriter) bool {
	if h.limiter.Allow() {
		return true
	}

	rateLimit := strconv.Itoa(int(h.limiter.Limit()))
	w.Header().Set("Retry-After", "1")
	w.Header().Set("X-RateLimit-Limit", rateLimit)
	http.Error(w, "too many requests", http.StatusTooManyRequests)
	return false
}

func createLimiter(limit int) *rate.Limiter {
	if limit <= 0 {
		return rate.NewLimiter(rate.Inf, 0)
	}
	return rate.NewLimiter(rate.Limit(limit), 1)
}
