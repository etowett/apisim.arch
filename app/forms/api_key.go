package forms

import "github.com/revel/revel"

type (
	ApiKey struct {
		UserID   int64
		Provider string
		Name     string
		Username string
		DlrURL   string
	}

	ApiKeyDlr struct {
		DlrURL string
	}
)

func (form *ApiKey) Validate(v *revel.Validation) {
	v.Required(form.UserID).Message("UserID required")
	v.Required(form.Provider).Message("Provider required")
	v.Required(form.Name).Message("Name required")
	v.Required(form.Username).Message("Username required")
}

func (form *ApiKeyDlr) Validate(v *revel.Validation) {
	v.Required(form.DlrURL).Message("DlrURL required")
}
