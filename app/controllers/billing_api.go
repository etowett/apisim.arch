package controllers

import (
	"apisim/app/db"
	"apisim/app/entities"
	"apisim/app/forms"
	"apisim/app/models"
	"apisim/app/webutils"
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/revel/revel"
)

type BillingAPI struct {
	App
}

func (c BillingAPI) Mpesa() revel.Result {
	var status int
	mpesaForm := forms.MpesaTopup{}
	err := c.Params.BindJSON(&mpesaForm)
	if err != nil {
		status = http.StatusBadRequest
		c.Response.SetStatus(status)
		return c.RenderJSON(entities.Response{
			Success: false,
			Status:  status,
			Message: fmt.Sprintf("failed to decode provided form: %v", err),
		})
	}

	c.Log.Infof("mpesa topup request: =[%+v]", mpesaForm)

	v := c.Validation
	mpesaForm.Validate(v)
	var errors []string
	for _, e := range v.Errors {
		errors = append(errors, strings.TrimSuffix(e.Message, "\n"))
	}

	if v.HasErrors() {
		status = http.StatusBadRequest
		c.Response.SetStatus(status)
		return c.RenderJSON(entities.Response{
			Success: false,
			Status:  status,
			Message: fmt.Sprintf("invalid form provided: =[%v]", strings.Join(errors, ",")),
		})
	}

	newUser := &models.User{}
	theUser, err := newUser.ByUsername(c.Request.Context(), db.DB(), mpesaForm.Username)
	if err != nil && err != sql.ErrNoRows {
		c.Log.Errorf("error getting user by username for topup: %v", err)
		status = http.StatusBadRequest
		c.Response.SetStatus(status)
		return c.RenderJSON(entities.Response{
			Success: false,
			Status:  status,
			Message: fmt.Sprintf("invalid username [%v] provided", mpesaForm.Username),
		})
	}

	if theUser.ID < 1 {
		status = http.StatusBadRequest
		c.Response.SetStatus(status)
		return c.RenderJSON(entities.Response{
			Success: false,
			Status:  status,
			Message: fmt.Sprintf("invalid username [%v] provided", mpesaForm.Username),
		})
	}

	newTrans := &models.Transaction{}
	lastTrans, err := newTrans.LastTransaction(c.Request.Context(), db.DB(), theUser.ID)
	if err != nil {
		if err != sql.ErrNoRows {
			c.Log.Errorf("error getting balance for mpesa top up: %v", err)
			status = http.StatusInternalServerError
			c.Response.SetStatus(status)
			return c.RenderJSON(entities.Response{
				Success: false,
				Status:  status,
				Message: "internal server error occured",
			})
		}
	}

	trans := models.Transaction{
		Amount:   mpesaForm.Amount,
		Currency: "KES",
		Code:     mpesaForm.MpesaCode,
		Type:     "mpesa_topup",
		UserID:   theUser.ID,
		Balance:  lastTrans.Balance + mpesaForm.Amount,
	}

	err = trans.Save(c.Request.Context(), db.DB())
	if err != nil {
		c.Log.Errorf("error saving for mpesa top up: %v", err)
		status = http.StatusInternalServerError
		c.Response.SetStatus(status)
		return c.RenderJSON(entities.Response{
			Success: false,
			Status:  status,
			Message: "internal server error occured",
		})
	}

	status = http.StatusCreated
	c.Response.SetStatus(status)
	return c.RenderJSON(entities.Response{
		Success: true,
		Status:  status,
		Message: "Topup Successful",
		Data:    trans,
	})
}

func (c BillingAPI) ForUser(id int64) revel.Result {
	status := http.StatusOK
	paginationFilter, err := webutils.FilterFromQuery(c.Params)
	if err != nil {
		c.Log.Errorf("could not filter from params: %v", err)
		status = http.StatusBadRequest
		c.Response.Status = status
		return c.RenderJSON(entities.Response{
			Success: false,
			Status:  status,
			Message: "Failed to parse page filters",
		})
	}

	newTrans := &models.Transaction{}
	data, err := newTrans.AllForUser(c.Request.Context(), db.DB(), id, paginationFilter)
	if err != nil {
		c.Log.Errorf("could not get transactions for user %v: %v", id, err)
		status = http.StatusInternalServerError
		c.Response.Status = status
		return c.RenderJSON(entities.Response{
			Success: false,
			Status:  status,
			Message: "Could not get transactions",
		})
	}

	recordsCount, err := newTrans.Count(c.Request.Context(), db.DB(), id, paginationFilter)
	if err != nil {
		c.Log.Errorf("could not get transactions count: %v", err)
		status = http.StatusInternalServerError
		c.Response.Status = status
		return c.RenderJSON(entities.Response{
			Success: false,
			Status:  status,
			Message: "Could not get transactions count",
		})
	}

	return c.RenderJSON(entities.Response{
		Success: true,
		Status:  status,
		Data: map[string]interface{}{
			"transactions": data,
			"pagination":   models.NewPagination(recordsCount, paginationFilter.Page, paginationFilter.Per),
		},
	})
}
