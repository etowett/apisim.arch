package entities

type (
	// ATResponse struct for at response
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
)
