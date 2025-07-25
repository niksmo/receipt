package service

import (
	"context"
	"fmt"
	"net/smtp"
)

type SMTPTransmitter struct {
	login, password, host, port string
}

func NewSMTPTransmitter(login, password, host, port string) SMTPTransmitter {
	return SMTPTransmitter{login, password, host, port}
}

func (t SMTPTransmitter) Send(ctx context.Context, to string, sub string, payload []byte) error {
	const op = "SMTPTransmitter.Send"

	msgTo := fmt.Sprintf("To: %s\r\n", to)
	msgSub := fmt.Sprintf("Subject: %s\r\n", sub)

	msg := []byte(msgTo + msgSub)
	msg = append(msg, payload...)
	msg = append(msg, []byte("\r\n")...)

	err := t.send(ctx, to, msg)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (t SMTPTransmitter) send(ctx context.Context, to string, msg []byte) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	doneStream := make(chan error)
	defer close(doneStream)

	go t.sendMail(doneStream, to, msg)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-doneStream:
		return err
	}
}

func (t SMTPTransmitter) sendMail(done chan error, to string, msg []byte) {
	auth := smtp.PlainAuth("", t.login, t.password, t.host)
	done <- smtp.SendMail(
		t.addr(), auth, "", []string{to}, msg,
	)
}

func (t SMTPTransmitter) addr() string {
	return t.host + ":" + t.port
}
