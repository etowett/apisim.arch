package forms

import (
	"regexp"

	"github.com/revel/revel"
)

var userRegex = regexp.MustCompile("^\\w*$")
var emailRegex = regexp.MustCompile(`[a-zA-Z0-9_\-]+@[a-zA-Z0-9_\-]+\.[a-zA-Z0-9_\-]+[a-zA-Z0-9]+$`)

type (
	User struct {
		Username             string `json:"username"`
		FirstName            string `json:"first_name"`
		LastName             string `json:"last_name"`
		Email                string `json:"email"`
		Password             string `json:"password"`
		PasswordConfirmation string `json:"password_confirmation"`
	}
)

func (user *User) Validate(v *revel.Validation) {
	ValidateUsername(v, user.Username).Key("user.Username")
	ValidatePassword(v, user.Password).Key("user.Password")
	ValidateFirstName(v, user.FirstName).Key("user.FirstName")
	ValidateLastName(v, user.LastName).Key("user.LastName")
	ValidateEmail(v, user.Email).Key("user.Email")

	v.Required(user.Password == user.PasswordConfirmation).
		Message("The passwords do not match")
}

func ValidateUsername(v *revel.Validation, username string) *revel.ValidationResult {
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

func ValidateFirstName(v *revel.Validation, firstName string) *revel.ValidationResult {
	result := v.Required(firstName).Message("FirstName is required")
	if !result.Ok {
		return result
	}

	result = v.MinSize(firstName, 3).Message("FirstName must exceed 2 characters")
	if !result.Ok {
		return result
	}

	result = v.MaxSize(firstName, 50).Message("FirstName cannot exceed 200 characters")

	return result
}
func ValidateLastName(v *revel.Validation, lastName string) *revel.ValidationResult {
	result := v.Required(lastName).Message("LastName is required")
	if !result.Ok {
		return result
	}

	result = v.MinSize(lastName, 3).Message("LastName must exceed 2 characters")
	if !result.Ok {
		return result
	}

	result = v.MaxSize(lastName, 50).Message("LastName cannot exceed 200 characters")

	return result
}

func ValidateEmail(v *revel.Validation, email string) *revel.ValidationResult {
	result := v.Required(email).Message("Email is required")
	if !result.Ok {
		return result
	}

	result = v.MinSize(email, 3).Message("Email must exceed 2 characters")
	if !result.Ok {
		return result
	}

	result = v.Match(email, emailRegex).Message("Email must be valid")
	if !result.Ok {
		return result
	}

	result = v.MaxSize(email, 50).Message("Email cannot exceed 200 characters")

	return result
}

func ValidatePassword(v *revel.Validation, password string) *revel.ValidationResult {
	return v.Check(password,
		revel.Required{},
		revel.MaxSize{100},
		revel.MinSize{5},
	)
}
