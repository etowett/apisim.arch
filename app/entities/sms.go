package entities

import "time"

type (
	ATResponse struct {
		Message    string         `json:"Message"`
		Recipients []*ATRecipient `json:"Recipients"`
	}

	ATRecipient struct {
		Number    string `json:"number"`
		Status    string `json:"status"`
		Cost      string `json:"cost"`
		MessageID string `json:"messageId"`
	}

	ProcessRequest struct {
		UserID     int64
		Meta       string
		Route      string
		Recipients []*ATRecipient
		SenderID   string
		Message    string
		Currency   string
		Cost       float64
		SentAt     time.Time
		StatusURL  string
	}
)
