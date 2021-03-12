package forms

import "github.com/revel/revel"

type (
	MpesaTopup struct {
		Username  string  `json:"username"`
		Amount    float64 `json:"amount"`
		Number    string  `json:"number"`
		TransTime string  `json:"trans_time"`
		MpesaCode string  `json:"mpesa_code"`
		KYCInfo   string  `json:"kyc_info"`
	}
)

func (form *MpesaTopup) Validate(v *revel.Validation) {
	v.Required(form.Username).Message("Username required")
	v.Required(form.Amount).Message("Amount required")
	v.Required(form.Number).Message("Number required")
	v.Required(form.TransTime).Message("Trans time required")
	v.Required(form.MpesaCode).Message("MpesaCode required")
	v.Required(form.KYCInfo).Message("KYCInfo required")
}
