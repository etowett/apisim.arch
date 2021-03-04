package controllers

import (
	"apisim/app/db"
	"apisim/app/forms"
	"apisim/app/models"
	"net/http"

	"github.com/revel/revel"
	"golang.org/x/crypto/bcrypt"
)

type UsersAPI struct {
	App
}

func (c *UsersAPI) Save() revel.Result {
	userForm := forms.User{}
	c.Params.BindJSON(&userForm)

	v := c.Validation
	userForm.Validate(v)
	if v.HasErrors() {
		c.Log.Errorf("Failed to validate given user form")
		result := response(v.Errors, "error save form", "failed")
		c.Response.ContentType = "application/json"
		c.Response.SetStatus(http.StatusBadRequest)
		return c.RenderJSON(result)
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(userForm.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Log.Errorf("error generating password hash: %v", err)
		result := response(err, "error generating password hash", "failed")
		c.Response.ContentType = "application/json"
		c.Response.SetStatus(http.StatusInternalServerError)
		return c.RenderJSON(result)
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
		result := response(err, "error save user", "failed")
		c.Response.ContentType = "application/json"
		c.Response.SetStatus(http.StatusInternalServerError)
		return c.RenderJSON(result)
	}

	c.Response.SetStatus(http.StatusCreated)
	c.Response.ContentType = "application/json"
	result := response(newUser, "insert data successfull", "success")
	return c.RenderJSON(result)
}
