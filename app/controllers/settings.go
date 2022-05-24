package controllers

import (
	"github.com/revel/revel"
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
