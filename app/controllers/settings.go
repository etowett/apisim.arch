package controllers

import (
	"apisim/app/db"
	"apisim/app/forms"
	"apisim/app/models"
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

func (c Settings) ApiKeySaveDlr(id int64, form *forms.ApiKeyDlr) revel.Result {
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

	return c.Redirect(ApiKeys.List)
}

func (c Settings) DeleteApiUrl(id int64) revel.Result {
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
