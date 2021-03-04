package controllers

import (
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/revel/revel"

	"apisim/app/entities"
	"apisim/app/forms"
	"apisim/app/helpers"
)

type SMSApi struct {
	App
}

func (c *SMSApi) inBlacklist(phoneNumber string) bool {
	// To check from db if inBlacklist
	return false
}

func (s *SMSApi) getMesageCost(
	message string,
	number string,
) float64 {
	return 1.0 * math.Ceil(float64(len(message))/160)
}

func (c SMSApi) SendToAT() revel.Result {

	smsRequest := &forms.ATForm{
		Username: c.Params.Get("username"),
		SenderID: c.Params.Get("from"),
		To:       c.Params.Get("to"),
		Message:  c.Params.Get("message"),
	}

	apiKey := c.Request.Header.Get("apikey")

	if len(apiKey) < 1 {
		return c.RenderText("Request is missing required HTTP header apikey.")
	}

	v := c.Validation
	smsRequest.Validate(v)

	if v.HasErrors() {
		retErrors := make([]string, 0)
		for _, theErr := range v.Errors {
			retErrors = append(retErrors, theErr.Message)
		}
		c.Response.SetStatus(http.StatusBadRequest)
		return c.RenderJSON(entities.ATResponse{
			Recipients: []*entities.ATRecipient{},
			Message:    strings.Join(retErrors, ","),
		})
	}

	if strings.ToLower(smsRequest.SenderID) == "testfail" {
		c.Response.SetStatus(http.StatusBadRequest)
		return c.RenderJSON(entities.ATResponse{
			Recipients: []*entities.ATRecipient{},
			Message:    fmt.Sprintf("Please ensure the originator [%v] is a shortCode or alphanumeric that is registered with us", smsRequest.SenderID),
		})
	}

	recipients := make([]*entities.ATRecipient, 0)
	validRecipients := make([]*entities.ATRecipient, 0)
	validCount := 0
	totalCost := 0.0
	for _, givenNumber := range strings.Split(smsRequest.To, ",") {
		var cost = 0.0
		var status = "Failed"
		var messageID = "None"

		validNum, err := helpers.GetValidPhone(givenNumber)
		if err != nil {
			status = "Invalid Phone Number"
		} else if c.inBlacklist(givenNumber) {
			status = "User In BlackList"
		} else {
			givenNumber = validNum
			status = "Success"
			cost = c.getMesageCost(smsRequest.Message, validNum)
			messageID = helpers.GetMD5Hash(time.Now().String() + givenNumber)
			validCount++
			totalCost += cost
			validRecipients = append(validRecipients, &entities.ATRecipient{
				Status:    status,
				Cost:      fmt.Sprintf("%.2f", cost),
				Number:    givenNumber,
				MessageID: messageID,
			})
		}

		recipients = append(recipients, &entities.ATRecipient{
			Status:    status,
			Cost:      fmt.Sprintf("%.2f", cost),
			Number:    givenNumber,
			MessageID: messageID,
		})
	}
	c.Log.Infof("AFT: (%d) - %s ", len(strings.Split(smsRequest.To, ",")), smsRequest.Message)

	retMessage := fmt.Sprintf(
		"Sent to %v/%v Total Cost: KES %v", validCount,
		len(recipients), totalCost,
	)

	// err = s.processMessage(c.Request.Context(), &entities.ProcessRequest{
	// 	UserID:     validKey.UserID,
	// 	Meta:       retMessage,
	// 	Recipients: recipients,
	// 	SenderID:   atForm.SenderID,
	// 	Message:    atForm.Message,
	// 	SentAt:     time.Now(),
	// })
	// if err != nil {
	// 	c.Log.Errorf("could not queue at message for process: %v", err)
	// }

	return c.RenderJSON(entities.ATResponse{
		Recipients: recipients,
		Message:    retMessage,
	})
}
