package controllers

import (
	"apisim/app/db"
	"apisim/app/entities"
	"apisim/app/forms"
	"apisim/app/helpers"
	"apisim/app/models"
	"apisim/app/webutils"
	"database/sql"
	"time"

	"github.com/revel/revel"
	null "gopkg.in/guregu/null.v4"
)

type Settings struct {
	App
}

func (c Settings) checkUser() revel.Result {
	if user := c.connected(); user == nil {
		c.Flash.Error("Please log in first")
		return c.Redirect(App.Index)
	}
	return nil
}

func (c Settings) Index() revel.Result {
	return c.Render()
}

func (c Settings) ApiKeyAdd() revel.Result {
	return c.Render()
}

func (c Settings) ApiKeySave(apiKey *forms.ApiKey) revel.Result {
	v := c.Validation
	apiKey.Validate(v)

	if v.HasErrors() {
		v.Keep()
		c.FlashParams()
		return c.Redirect(Settings.ApiKeyAdd)
	}

	theApiKey := &models.ApiKey{}
	existingApiKey, err := theApiKey.ByUserAndAccessID(c.Request.Context(), db.DB(), apiKey.Username, apiKey.UserID)
	if err != nil && err != sql.ErrNoRows {
		c.Log.Errorf("Failed theApiKey.ByUserAndAccessID =[%v]", err)
		c.Flash.Error("Internal server issue occured, please retry.")
		return c.Redirect(Settings.ApiKeyAdd)
	}

	if existingApiKey.ID > 1 {
		c.Flash.Error("You already have a key with name =[%v]", apiKey.Name)
		return c.Redirect(Settings.ApiKeyAdd)
	}

	accessSecret := helpers.GenerateApiKeySecret()
	accessSecretHash, err := helpers.HashApiKeySecret(accessSecret)
	if err != nil {
		c.Log.Errorf("Failed to hash api key secret for user=[%v]", apiKey.UserID)
		c.Flash.Error("Could not get secret, internal server issue.")
		return c.Redirect(Settings.ApiKeyAdd)
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
		c.Log.Errorf("Could not save apiKey: %v", err)
		c.Validation.Keep()
		c.Flash.Error("Could not save, internal server issue.")
		c.FlashParams()
		return c.Redirect(Settings.ApiKeyAdd)
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

	c.Flash.Success("ApiKey created - " + newApiKey.Name)
	c.Session["api-key-secret"] = accessSecret
	return c.Redirect(Settings.ApiKeyDetails, newApiKey.ID)
}

func (c Settings) ApiKeySaveDlr(id int64, form *forms.ApiKeyDlr) revel.Result {
	v := c.Validation
	form.Validate(v)

	if v.HasErrors() {
		v.Keep()
		c.FlashParams()
		return c.Redirect(Settings.ApiKeyDetails, id)
	}

	newApiKey := &models.ApiKey{}
	apiKey, err := newApiKey.ByID(c.Request.Context(), db.DB(), id)
	if err != nil {
		c.Log.Errorf("Could not get apiKey by id %v: %v", id, err)
		c.Validation.Keep()
		c.Flash.Error("Could not save, internal server issue.")
		c.FlashParams()
		return c.Redirect(Settings.ApiKeyDetails, id)
	}

	apiKey.DlrURL = form.DlrURL
	apiKey.UpdatedAt = null.TimeFrom(time.Now())

	err = apiKey.Save(c.Request.Context(), db.DB())
	if err != nil {
		c.Log.Errorf("Failed to save when updating apiKey: %v", err)
		c.Validation.Keep()
		c.Flash.Error("Could not save, internal server issue.")
		c.FlashParams()
		return c.Redirect(Settings.ApiKeyDetails, id)
	}

	err = c.clearCachedApiKey(apiKey.Provider, apiKey.AccessID)
	if err != nil {
		c.Log.Errorf("Failed to clearCachedApiKey: %v", err)
	}

	return c.Redirect(Settings.ApiKeyDetails, id)
}

func (c Settings) ApiKeys() revel.Result {
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

func (c Settings) ApiKeyDetails(id int64) revel.Result {
	var result entities.Response
	var theSecret string
	newApiKey := &models.ApiKey{}
	apiKey, err := newApiKey.ByID(c.Request.Context(), db.DB(), id)
	if err != nil {
		c.Log.Errorf("could not get apikey with id %v: %v", id, err)
		result = entities.Response{
			Success: false,
			Message: "Could not get apikey details",
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

func (c Settings) DeleteApiKey(id int64) revel.Result {
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

	return c.Redirect(Settings.ApiKeys)
}

func (c Settings) DeleteApiUrl(id int64) revel.Result {
	newApiKey := &models.ApiKey{}
	apiKey, err := newApiKey.ByID(c.Request.Context(), db.DB(), id)
	if err != nil {
		c.Log.Errorf("Failed to get for delete apiKey %v: %v", id, err)
		c.Validation.Keep()
		c.Flash.Error("Could not delete, internal server issue.")
		c.FlashParams()
		return c.Redirect(Settings.ApiKeyDetails, id)
	}

	apiKey.DlrURL = ""
	apiKey.UpdatedAt = null.TimeFrom(time.Now())

	err = apiKey.Save(c.Request.Context(), db.DB())
	if err != nil {
		c.Log.Errorf("Failed to save when updating apiKey dlr: %v", err)
		c.Validation.Keep()
		c.Flash.Error("Could not delete, internal server issue.")
		c.FlashParams()
		return c.Redirect(Settings.ApiKeyDetails, id)
	}

	err = c.clearCachedApiKey(apiKey.Provider, apiKey.AccessID)
	if err != nil {
		c.Log.Errorf("Failed to clearCachedApiKey: %v", err)
	}

	return c.Redirect(Settings.ApiKeyDetails, id)
}
