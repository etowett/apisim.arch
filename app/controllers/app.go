package controllers

import (
	"apisim/app"
	"apisim/app/db"
	"apisim/app/entities"
	"apisim/app/models"
	"encoding/json"
	"fmt"
	"time"

	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	loggedInUser := c.connected()
	if loggedInUser != nil {
		return c.Redirect(App.Dash)
	}
	return c.Render()
}

func (c App) Dash() revel.Result {
	loggedInUser := c.connected()
	if loggedInUser == nil {
		return c.Redirect(App.Index)
	}
	return c.Render()
}

func (c App) Health() revel.Result {
	return c.RenderJSON(map[string]interface{}{
		"success":     true,
		"status":      "Ok",
		"server_time": time.Now(),
		"version":     app.AppVersion,
		"build_time":  app.BuildTime,
	})
}

func (c App) getUserFromUsername(username string) *models.User {
	user := &models.User{}
	c.Session.GetInto("user", user, false)
	if user.Username == username {
		return user
	}

	newUser := &models.User{}
	foundUser, err := newUser.ByUsername(c.Request.Context(), db.DB(), username)
	if err != nil {
		c.Log.Errorf("could not get user by username: %v", err)
		return nil
	}

	c.Session["user"] = foundUser
	return foundUser
}

func (c App) connected() *models.User {
	if c.ViewArgs["user"] != nil {
		return c.ViewArgs["user"].(*models.User)
	}
	if username, ok := c.Session["username"]; ok {
		return c.getUserFromUsername(username.(string))
	}
	return nil
}

func (c App) AddUser() revel.Result {
	if user := c.connected(); user != nil {
		c.ViewArgs["user"] = user
	}
	return nil
}

func (c App) cacheApiKey(
	net string,
	accountID string,
	cachedApiKey *entities.CachedApiKey,
) error {

	b, err := json.Marshal(cachedApiKey)
	if err != nil {
		return fmt.Errorf("failed to marshal for cache api key for user=[%v]: %v", cachedApiKey.UserID, err)
	}

	_, err = redisManager.Set(c.generateCacheKey(net, accountID), b)
	if err != nil {
		return fmt.Errorf("failed to cache api key for user=[%v]: %v", cachedApiKey.UserID, err)
	}

	return nil
}

func (c App) clearCachedApiKey(
	net string,
	accountID string,
) error {
	if _, err := redisManager.Del(c.generateCacheKey(net, accountID)); err != nil {
		return fmt.Errorf("failed to clear cached api key for accountid=[%v]: %v", accountID, err)
	}
	return nil
}

func (c App) generateCacheKey(net, accountID string) string {
	return "apisim:" + net + ":apikey:" + accountID
}
