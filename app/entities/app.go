package entities

type (
	Response struct {
		Data    interface{} `json:"data"`
		Message string      `json:"message"`
		Success bool        `json:"success"`
		Status  int         `json:"status"`
	}
)
