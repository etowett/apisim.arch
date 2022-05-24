package controllers

import (
	"apisim/app/db"
	"apisim/app/entities"
	"apisim/app/forms"
	"apisim/app/helpers"
	"apisim/app/models"
	"apisim/app/webutils"
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/revel/revel"
	null "gopkg.in/guregu/null.v4"
)

type (
	ApiKeys struct {
		App
	}
)

func (c ApiKeys) Add() revel.Result {
	return c.Render()
}

func (c ApiKeys) List() revel.Result {
	var result entities.Response
	paginationFilter, err := webutils.FilterFromQuery(c.Params)
	if err != nil {
		c.Log.Errorf("could not filter from params: %v", err)
		result = entities.Response{
			Success: false,
			Message: "Failed to parse page filters",
		}
		return c.Render(result)
	}

	newApiKey := &models.ApiKey{}
	data, err := newApiKey.All(c.Request.Context(), db.DB(), paginationFilter)
	if err != nil {
		c.Log.Errorf("could not get apikeys: %v", err)
		result = entities.Response{
			Success: false,
			Message: "Could not get apikeys",
		}
		return c.Render(result)
	}

	recordsCount, err := newApiKey.Count(c.Request.Context(), db.DB(), paginationFilter)
	if err != nil {
		c.Log.Errorf("could not get apikeys count: %v", err)
		result = entities.Response{
			Success: false,
			Message: "Could not get apikeys count",
		}
		return c.Render(result)
	}

	result = entities.Response{
		Success: true,
		Data: map[string]interface{}{
			"ApiKeys":    data,
			"Pagination": models.NewPagination(recordsCount, paginationFilter.Page, paginationFilter.Per),
		},
	}
	return c.Render(result)
}

func (c ApiKeys) APIDetails(id int64) revel.Result {
	apiKey, status, appErr, err := c.getApiDetails(c.Request.Context(), id)

	if err != nil {
		if appErr != nil {
			c.Log.Errorf("failed api getApiDetails - %+v", appErr)
		}
		return c.RenderJSON(entities.Response{
			Message: err.Error(),
			Status:  status,
			Success: false,
			Data:    apiKey,
		})
	}

	return c.RenderJSON(entities.Response{
		Data:    apiKey,
		Status:  status,
		Message: "Fetch success",
		Success: true,
	})
}

func (c ApiKeys) Delete(id int64) revel.Result {
	newApiKey := &models.ApiKey{}
	theApiKey, err := newApiKey.ByID(c.Request.Context(), db.DB(), id)
	if err != nil {
		c.Log.Errorf("error newApiKey by id %v: %v", id, err)
	}

	_, err = redisManager.Del(c.generateCacheKey(theApiKey.Provider, theApiKey.AccessID))
	if err != nil {
		c.Log.Errorf("failed to delete apikey in redis for id %v: %v", id, err)
	}

	err = theApiKey.Delete(c.Request.Context(), db.DB(), id)
	if err != nil {
		c.Log.Errorf("error newApiKey delete: %v", err)
	}

	return c.Redirect(ApiKeys.List)
}

func (c ApiKeys) DeleteUrl(id int64) revel.Result {
	newApiKey := &models.ApiKey{}
	apiKey, err := newApiKey.ByID(c.Request.Context(), db.DB(), id)
	if err != nil {
		c.Log.Errorf("Failed to get for delete apiKey %v: %v", id, err)
		c.Validation.Keep()
		c.Flash.Error("Could not delete, internal server issue.")
		c.FlashParams()
		return c.Redirect(ApiKeys.Details, id)
	}

	apiKey.DlrURL = ""
	apiKey.UpdatedAt = null.TimeFrom(time.Now())

	err = apiKey.Save(c.Request.Context(), db.DB())
	if err != nil {
		c.Log.Errorf("Failed to save when updating apiKey dlr: %v", err)
		c.Validation.Keep()
		c.Flash.Error("Could not delete, internal server issue.")
		c.FlashParams()
		return c.Redirect(ApiKeys.Details, id)
	}

	err = c.clearCachedApiKey(apiKey.Provider, apiKey.AccessID)
	if err != nil {
		c.Log.Errorf("Failed to clearCachedApiKey: %v", err)
	}

	return c.Redirect(ApiKeys.Details, id)
}

func (c ApiKeys) getApiDetails(
	ctx context.Context,
	id int64,
) (*models.ApiKey, int, error, error) {
	newApiKey := &models.ApiKey{}
	apiKey, err := newApiKey.ByID(ctx, db.DB(), id)
	if err != nil {
		return apiKey,
			http.StatusInternalServerError,
			fmt.Errorf("could not get apikey with id %v: %v", id, err),
			fmt.Errorf("Could not get apikey details")
	}

	return apiKey, http.StatusOK, nil, nil
}

func (c ApiKeys) Details(id int64) revel.Result {
	var result entities.Response
	var theSecret string

	apiKey, _, err, userErr := c.getApiDetails(c.Request.Context(), id)
	if userErr != nil {
		if err != nil {
			c.Log.Errorf("failed to get apikey: %v", err)
		}
		result = entities.Response{
			Success: false,
			Message: userErr.Error(),
		}
		return c.Render(result)
	}

	if keySecret, ok := c.Session["api-key-secret"]; ok {
		theSecret = keySecret.(string)
		delete(c.Session, "api-key-secret")
	}

	result = entities.Response{
		Success: true,
		Data: map[string]interface{}{
			"ApiKey": apiKey,
			"Secret": theSecret,
		},
	}
	return c.Render(result)
}

func (c ApiKeys) Save(apiKey *forms.ApiKey) revel.Result {
	v := c.Validation
	apiKey.Validate(v)

	ctx := c.Request.Context()

	c.Log.Infof("request: %+v", c.Request.Form)

	if v.HasErrors() {
		v.Keep()
		c.FlashParams()
		return c.Redirect(ApiKeys.Add)
	}

	err, appErr, _, accessSecret, newApiKey := c.createApiKey(ctx, apiKey)
	if err != nil {
		if appErr != nil {
			c.Log.Errorf("failed createApiKey - %+v", appErr)
		}
		c.Validation.Keep()
		c.Flash.Error(err.Error())
		c.FlashParams()
		return c.Redirect(ApiKeys.Add)
	}

	c.Flash.Success("ApiKey created - " + newApiKey.Name)
	c.Session["api-key-secret"] = accessSecret
	return c.Redirect(ApiKeys.Details, newApiKey.ID)
}

func (c ApiKeys) ApiCreate() revel.Result {
	var status int
	apiKeyForm := forms.ApiKey{}

	c.Params.BindJSON(&apiKeyForm)

	ctx := c.Request.Context()
	v := c.Validation
	apiKeyForm.Validate(v)
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

	err, appErr, status, accessSecret, apiKey := c.createApiKey(ctx, &apiKeyForm)
	c.Response.SetStatus(status)
	if err != nil {
		if appErr != nil {
			c.Log.Errorf("failed createApiKey - %+v", appErr)
		}
		return c.RenderJSON(entities.Response{
			Message: err.Error(),
			Status:  status,
			Success: false,
		})
	}

	return c.RenderJSON(entities.Response{
		Data: map[string]interface{}{
			"apikey": apiKey,
			"secret": accessSecret,
		},
		Status:  status,
		Success: true,
	})
}

func (c ApiKeys) createApiKey(
	ctx context.Context,
	apiKey *forms.ApiKey,
) (error, error, int, string, *models.ApiKey) {
	theApiKey := &models.ApiKey{}
	existingApiKey, err := theApiKey.ByUserAndAccessID(ctx, db.DB(), apiKey.Username, apiKey.UserID)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("Internal server issue occured, please retry."),
			fmt.Errorf("Failed theApiKey.ByUserAndAccessID =[%v]", err),
			http.StatusBadRequest,
			"",
			theApiKey
	}

	if existingApiKey.ID > 1 {
		return fmt.Errorf("You already have a key with name =[%v]", apiKey.Name),
			nil,
			http.StatusConflict,
			"",
			theApiKey
	}

	accessSecret := helpers.GenerateApiKeySecret()
	accessSecretHash, err := helpers.HashApiKeySecret(accessSecret)
	if err != nil {
		return fmt.Errorf("Could not get secret, internal server error"),
			fmt.Errorf("Failed to hash api key secret for user=[%v]", apiKey.UserID),
			http.StatusInternalServerError,
			accessSecret,
			theApiKey
	}

	newApiKey := &models.ApiKey{
		UserID:           apiKey.UserID,
		Provider:         apiKey.Provider,
		Name:             apiKey.Name,
		AccessID:         apiKey.Username,
		AccessSecretHash: accessSecretHash,
		DlrURL:           apiKey.DlrURL,
	}
	err = newApiKey.Save(c.Request.Context(), db.DB())
	if err != nil {
		return fmt.Errorf("Could not save the api key, internal server error"),
			fmt.Errorf("could not save apiKey: %v", err),
			http.StatusInternalServerError,
			accessSecret,
			newApiKey
	}

	cachedApiKey := &entities.CachedApiKey{
		UserID:            apiKey.UserID,
		AccountSecretHash: accessSecretHash,
		DlrURL:            apiKey.DlrURL,
	}

	err = c.cacheApiKey(apiKey.Provider, apiKey.Username, cachedApiKey)
	if err != nil {
		c.Log.Errorf("could not cache api key: %v", err)
	}

	return nil,
		nil,
		http.StatusCreated,
		accessSecret,
		newApiKey
}

func (c ApiKeys) SaveDlr(id int64, form *forms.ApiKeyDlr) revel.Result {
	v := c.Validation
	form.Validate(v)

	if v.HasErrors() {
		v.Keep()
		c.FlashParams()
		return c.Redirect(ApiKeys.Details, id)
	}

	newApiKey := &models.ApiKey{}
	apiKey, err := newApiKey.ByID(c.Request.Context(), db.DB(), id)
	if err != nil {
		c.Log.Errorf("Could not get apiKey by id %v: %v", id, err)
		c.Validation.Keep()
		c.Flash.Error("Could not save, internal server issue.")
		c.FlashParams()
		return c.Redirect(ApiKeys.Details, id)
	}

	apiKey.DlrURL = form.DlrURL
	apiKey.UpdatedAt = null.TimeFrom(time.Now())

	err = apiKey.Save(c.Request.Context(), db.DB())
	if err != nil {
		c.Log.Errorf("Failed to save when updating apiKey: %v", err)
		c.Validation.Keep()
		c.Flash.Error("Could not save, internal server issue.")
		c.FlashParams()
		return c.Redirect(ApiKeys.Details, id)
	}

	err = c.clearCachedApiKey(apiKey.Provider, apiKey.AccessID)
	if err != nil {
		c.Log.Errorf("Failed to clearCachedApiKey: %v", err)
	}

	return c.Redirect(ApiKeys.Details, id)
}
