package controllers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/revel/revel"

	"apisim/app/db"
	"apisim/app/entities"
	"apisim/app/forms"
	"apisim/app/helpers"
	"apisim/app/jobs/sms_jobs"
	"apisim/app/models"
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

func (c *SMSApi) validateApiKey(
	ctx context.Context,
	net string,
	accountID string,
	accountSecret string,
) (*entities.CachedApiKey, error) {

	cachedApiKey := &entities.CachedApiKey{}

	doCache := false
	val, err := redisManager.GetString(c.generateCacheKey(net, accountID))
	if err != nil {
		if redisManager.IsErrNil(err) {
			doCache = true
			newApiKey := &models.ApiKey{}
			apiKey, err := newApiKey.ByAccountID(ctx, db.DB(), accountID)
			if err != nil {
				if err == sql.ErrNoRows {
					return cachedApiKey, fmt.Errorf("Invalid api credentials provided")
				}

				return cachedApiKey, fmt.Errorf("Failed to get api key when validating api key")
			}

			cachedApiKey.UserID = apiKey.UserID
			cachedApiKey.AccountSecretHash = apiKey.AccessSecretHash
		} else {
			return cachedApiKey, fmt.Errorf("Failed to retrieve api key from cache")
		}
	} else {
		err = json.Unmarshal([]byte(val), cachedApiKey)
		if err != nil {
			return cachedApiKey, fmt.Errorf("Failed to marshal cached api key: %v", err)
		}
	}

	if !helpers.CheckPasswordHash(accountSecret, cachedApiKey.AccountSecretHash) {
		return cachedApiKey, fmt.Errorf("Invalid api credentials provided")
	}

	if doCache {
		err := c.cacheApiKey(net, accountID, cachedApiKey)
		if err != nil {
			c.Log.Errorf("could not cache api key: %v", err)
		}
	}

	return cachedApiKey, nil
}

func (c *SMSApi) SendToAT() revel.Result {

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

	if len(smsRequest.Username) < 1 {
		return c.RenderText("Must have username in your request.")
	}

	validKey, err := c.validateApiKey(c.Request.Context(), "at", smsRequest.Username, apiKey)
	if err != nil {
		return c.RenderText("The supplied authentication is invalid.")
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

	err = c.processMessage(c.Request.Context(), &entities.ProcessRequest{
		UserID:     validKey.UserID,
		Meta:       retMessage,
		Recipients: recipients,
		SenderID:   smsRequest.SenderID,
		Message:    smsRequest.Message,
		SentAt:     time.Now(),
	})
	if err != nil {
		c.Log.Errorf("could not queue at message for process: %v", err)
		c.Response.SetStatus(http.StatusInternalServerError)
		return c.RenderJSON(entities.ATResponse{
			Recipients: []*entities.ATRecipient{},
			Message:    "Server error. Could not process your message",
		})
	}

	return c.RenderJSON(entities.ATResponse{
		Recipients: recipients,
		Message:    retMessage,
	})
}

func (c *SMSApi) processMessage(
	ctx context.Context,
	smsRequest *entities.ProcessRequest,
) error {
	_, err := jobEnqueuer.Enqueue(ctx, sms_jobs.NewSendSMSJob(smsRequest))
	return err
	return nil
}
