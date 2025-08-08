package adapter

type Sender struct {
	Email string `json:"email"`
}

type SendTo struct {
	Email string `json:"email"`
}

type SendEmail struct {
	Sender      Sender   `json:"sender"`
	To          []SendTo `json:"to"`
	Subject     string   `json:"subject"`
	TextContent string   `json:"textContent"`
}

type MessageCreated struct {
	MessageID string `json:"messageId"`
}

type BadRequest struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
