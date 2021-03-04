package controllers

import "github.com/revel/revel"

type Outbox struct {
	App
}

func (c Outbox) All() revel.Result {
	return c.Render()
}

func (c Outbox) Get(id int64) revel.Result {
	return c.Render()
}
