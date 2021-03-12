package controllers

import (
	"apisim/app/db"
	"apisim/app/entities"
	"apisim/app/models"
	"apisim/app/webutils"
	"bytes"
	"fmt"
	"time"

	"github.com/revel/revel"
)

type Outbox struct {
	App
}

func (c Outbox) All() revel.Result {
	loggedInUser := c.connected()
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

	newMessage := &models.Message{}
	data, err := newMessage.AllForUser(c.Request.Context(), db.DB(), loggedInUser.ID, paginationFilter)
	if err != nil {
		c.Log.Errorf("could not get messages: %v", err)
		result = entities.Response{
			Success: false,
			Message: "Could not get messages",
		}
		return c.Render(result)
	}

	recordsCount, err := newMessage.Count(c.Request.Context(), db.DB(), loggedInUser.ID, paginationFilter)
	if err != nil {
		c.Log.Errorf("could not get messages count: %v", err)
		result = entities.Response{
			Success: false,
			Message: "Could not get messages count",
		}
		return c.Render(result)
	}

	result = entities.Response{
		Success: true,
		Data: map[string]interface{}{
			"Messages":   data,
			"Pagination": models.NewPagination(recordsCount, paginationFilter.Page, paginationFilter.Per),
		},
	}
	return c.Render(result)
}

func (c Outbox) Get(id int64) revel.Result {
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

	newMessage := &models.Message{}
	message, err := newMessage.ByID(c.Request.Context(), db.DB(), id)
	if err != nil {
		c.Log.Errorf("could not get message with id %v: %v", id, err)
		result = entities.Response{
			Success: false,
			Message: "Could not get message details",
		}
		return c.Render(result)
	}

	newRec := &models.Recipient{}
	data, err := newRec.ForMessage(c.Request.Context(), db.DB(), id, paginationFilter)
	if err != nil {
		c.Log.Errorf("could not get recipients for message %v: %v", id, err)
		result = entities.Response{
			Success: false,
			Message: "Could not get recipients",
		}
		return c.Render(result)
	}

	recordsCount, err := newRec.CountForMessage(c.Request.Context(), db.DB(), id, paginationFilter)
	if err != nil {
		c.Log.Errorf("could not get recipients count: %v", err)
		result = entities.Response{
			Success: false,
			Message: "Could not get recipients count",
		}
		return c.Render(result)
	}

	result = entities.Response{
		Success: true,
		Data: map[string]interface{}{
			"Message":    message,
			"Recipients": data,
			"Pagination": models.NewPagination(recordsCount, paginationFilter.Page, paginationFilter.Per),
		},
	}
	return c.Render(result)
}

func (c Outbox) ExportAll() revel.Result {
	loggedInUser := c.connected()
	c.Log.Infof("loggedInUser: %v", loggedInUser)

	newMessage := &models.Message{}
	data, err := newMessage.AllForUser(c.Request.Context(), db.DB(), 1, &models.Filter{})
	// data, err := newMessage.AllForUser(c.Request.Context(), db.DB(), loggedInUser.ID, &models.Filter{})
	if err != nil {
		c.Log.Errorf("could not get messages for export: %v", err)
		return nil
	}

	b, err := csvCreator.CreateMessagesCSV(data)
	if err != nil {
		c.Log.Errorf("Failed to create messages csv when exporting messages: %v", err)
		return nil
	}

	return c.RenderBinary(
		bytes.NewReader(b),
		fmt.Sprintf("my_messages_%s.csv", time.Now().Format("20060102150405")),
		revel.Attachment,
		time.Now(),
	)
}

func (c Outbox) ExportRecipients(id int64) revel.Result {
	newRec := &models.Recipient{}
	data, err := newRec.ForMessage(c.Request.Context(), db.DB(), id, &models.Filter{})
	if err != nil {
		c.Log.Errorf("could not get recipients for export %v: %v", id, err)
		return nil
	}

	b, err := csvCreator.CreateRecipentsCSV(data)
	if err != nil {
		c.Log.Errorf("Failed to create recipients csv when exporting: %v", err)
		return nil
	}

	return c.RenderBinary(
		bytes.NewReader(b),
		fmt.Sprintf("message_recipients_%s.csv", time.Now().Format("20060102150405")),
		revel.Attachment,
		time.Now(),
	)
}
