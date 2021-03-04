package forms

import "github.com/revel/revel"

type (
	Login struct {
		Username string
		Password string
		Remember bool
	}
)

func (login *Login) Validate(v *revel.Validation) {
	ValidateLoginUsername(v, login.Username).Key("login.Username")
	ValidateLoginPassword(v, login.Password).Key("login.Password")
}

func ValidateLoginUsername(v *revel.Validation, username string) *revel.ValidationResult {
	result := v.Required(username).Message("Username is required")
	if !result.Ok {
		return result
	}

	result = v.MinSize(username, 3).Message("Username must exceed 2 characters")
	if !result.Ok {
		return result
	}

	result = v.Match(username, userRegex).Message("Username must have valid characters")
	if !result.Ok {
		return result
	}

	result = v.MaxSize(username, 50).Message("Username cannot exceed 200 characters")

	return result
}

func ValidateLoginPassword(v *revel.Validation, password string) *revel.ValidationResult {
	return v.Check(password,
		revel.Required{},
		revel.MaxSize{100},
		revel.MinSize{3},
	)
}
