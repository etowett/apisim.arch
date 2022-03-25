package controllers

import (
	"apisim/app"
	"apisim/app/db"
	"apisim/app/entities"
	"apisim/app/models"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
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

func (c App) bytesToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func (c App) Health() revel.Result {
	ctx := c.Request.Context()

	hostName, err := os.Hostname()
	if err != nil {
		c.Log.Errorf("could not get hostname: %v", err)
	}

	dbUp := true
	err = db.DB().Ping()
	if err != nil {
		dbUp = false
		c.Log.Errorf("db ping failed because err=[%v]", err)
	}

	hello, err := (&models.Checker{}).One(ctx, db.DB())
	if err != nil {
		c.Log.Errorf("could not db select 1: %v", err)
	}

	version, err := (&models.Checker{}).Version(ctx, db.DB())
	if err != nil {
		c.Log.Errorf("could not db version: %v", err)
	}

	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)

	currentTime := time.Now()
	tZone, offset := currentTime.Zone()

	return c.RenderJSON(map[string]interface{}{
		"success": true,
		"time": map[string]interface{}{
			"now":      currentTime,
			"timezone": tZone,
			"offset":   offset,
		},
		"version":    app.AppVersion,
		"build_time": app.BuildTime,
		"db": map[string]interface{}{
			"type":    "postgres",
			"up":      dbUp,
			"hello":   hello,
			"version": version,
		},

		"server": map[string]interface{}{
			"hostname":   hostName,
			"cpu":        runtime.NumCPU(),
			"goroutines": runtime.NumGoroutine(),
			"goarch":     runtime.GOARCH,
			"goos":       runtime.GOOS,
			"compiler":   runtime.Compiler,
			"memory": map[string]interface{}{
				"alloc":       fmt.Sprintf("%v MB", c.bytesToMb(memStats.Alloc)),
				"total_alloc": fmt.Sprintf("%v MB", c.bytesToMb(memStats.TotalAlloc)),
				"sys":         fmt.Sprintf("%v MB", c.bytesToMb(memStats.Sys)),
				"num_gc":      memStats.NumGC,
			},
		},
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
