package forms

import "github.com/revel/revel"

type (
	ATForm struct {
		Username string `form:"username"`
		SenderID string `form:"from"`
		To       string `form:"to"`
		Message  string `form:"message"`
	}
)

func (form *ATForm) Validate(v *revel.Validation) {
	v.Required(form.Username).Message("Username - [username] is required")
	v.Required(form.SenderID).Message("SenderID - [from] is required")
	v.Required(form.To).Message("Destination - [to] is required")
	v.Required(form.Message).Message("Message - [message] is required")
}
