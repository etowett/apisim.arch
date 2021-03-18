package controllers

import (
	"apisim/app/db"
	"apisim/app/entities"
	"apisim/app/forms"
	"apisim/app/models"
	"net/http"
	"strings"

	"github.com/revel/revel"
	"golang.org/x/crypto/bcrypt"
)

type (
	UsersAPI struct {
		App
	}
)

func (c *UsersAPI) Save() revel.Result {
	var status int
	userForm := forms.User{}
	c.Params.BindJSON(&userForm)

	v := c.Validation
	userForm.Validate(v)
	if v.HasErrors() {
		retErrors := make([]string, 0)
		for _, theErr := range v.Errors {
			retErrors = append(retErrors, theErr.Message)
		}
		status = http.StatusBadRequest
		c.Response.SetStatus(status)
		return c.RenderJSON(entities.Response{
			Message: strings.Join(retErrors, ","),
			Status:  status,
			Success: false,
		})
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(userForm.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Log.Errorf("error generating password hash: %v", err)
		status = http.StatusInternalServerError
		c.Response.SetStatus(status)
		return c.RenderJSON(entities.Response{
			Message: "Encountered an error",
			Status:  status,
			Success: false,
		})
	}

	newUser := &models.User{
		Username:     userForm.Username,
		FirstName:    userForm.FirstName,
		LastName:     userForm.LastName,
		Email:        userForm.Email,
		PasswordHash: string(passwordHash[:]),
	}

	err = newUser.Save(c.Request.Context(), db.DB())
	if err != nil {
		c.Log.Errorf("error insert user: %v", err)
		status = http.StatusInternalServerError
		c.Response.SetStatus(status)
		return c.RenderJSON(entities.Response{
			Message: "Encountered an error saving request.",
			Status:  status,
			Success: false,
		})
	}

	status = http.StatusCreated
	c.Response.SetStatus(status)
	return c.RenderJSON(entities.Response{
		Data:    newUser,
		Status:  status,
		Success: true,
	})
}
