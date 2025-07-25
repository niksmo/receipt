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
	return SMTPTransmitter{}
}

func (t SMTPTransmitter) Send(ctx context.Context, to string, sub string, payload []byte) error {
	const op = "SMTPTransmitter.Send"

	to = fmt.Sprintf("To: %s\r\n", to)
	sub = fmt.Sprintf("Subject: %s\r\n", sub)

	msg := []byte(to + sub)
	msg = append(msg, payload...)
	msg = append(msg, []byte("\r\n")...)

	err := t.send(ctx, to, msg)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (t SMTPTransmitter) send(ctx context.Context, to string, msg []byte) error {
	auth := smtp.PlainAuth("",
		t.login, t.password, t.host)

	doneStream := make(chan error)
	defer close(doneStream)

	go func(done chan error) {
		err := smtp.SendMail(
			t.addr(), auth, "", []string{to}, msg,
		)
		done <- err
	}(doneStream)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-doneStream:
		return err
	}
}

func (t SMTPTransmitter) addr() string {
	return t.host + ":" + t.port
}
