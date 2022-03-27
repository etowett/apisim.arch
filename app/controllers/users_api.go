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

func (c UsersAPI) Save() revel.Result {
	var status int
	userForm := forms.User{}

	c.Log.Infof("got request: %+v", c.Request)

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

func (c UsersAPI) Login() revel.Result {
	var status int
	loginForm := forms.Login{}
	c.Params.BindJSON(&loginForm)

	v := c.Validation
	loginForm.Validate(v)
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

	user := c.getUserFromUsername(loginForm.Username)
	if user == nil {
		status = http.StatusBadRequest
		c.Response.SetStatus(status)
		return c.RenderJSON(entities.Response{
			Message: "Could not find user with that username",
			Status:  status,
			Success: false,
		})
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginForm.Password))
	if err != nil {
		status = http.StatusBadRequest
		c.Response.SetStatus(status)
		return c.RenderJSON(entities.Response{
			Message: "Invalid password provided",
			Status:  status,
			Success: false,
		})
	}

	status = http.StatusOK
	c.Response.SetStatus(status)
	return c.RenderJSON(entities.Response{
		Data:    user,
		Status:  status,
		Success: true,
	})
}
