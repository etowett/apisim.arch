package job_handlers

import (
	"apisim/app/db"
	"apisim/app/entities"
	"apisim/app/jobs"
	"apisim/app/jobs/sms_jobs"
	"apisim/app/models"
	"apisim/app/work"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/revel/revel"
)

type ProcessSMSJobHandler struct {
	jobEnqueuer work.JobEnqueuer
}

func NewProcessSMSJobHandler(
	jobEnqueuer work.JobEnqueuer,
) *ProcessSMSJobHandler {
	return &ProcessSMSJobHandler{
		jobEnqueuer: jobEnqueuer,
	}
}

func (h *ProcessSMSJobHandler) Job() jobs.Job {
	return &sms_jobs.SendSMSJob{}
}

func (h *ProcessSMSJobHandler) PerformJob(
	ctx context.Context,
	body string,
) error {
	var theJob sms_jobs.SendSMSJob
	err := json.Unmarshal([]byte(body), &theJob)
	if err != nil {
		revel.AppLog.Errorf("error unmarshal send sms task: %v", err)
		return nil
	}

	revel.AppLog.Infof("process send sms task: =[%+v]", theJob)

	req := theJob.Request
	message := models.Message{
		UserID:   req.UserID,
		SenderID: req.SenderID,
		Meta:     req.Message,
		Message:  req.Message,
		Cost:     req.Cost,
		Currency: req.Currency,
		SentAt:   req.SentAt,
	}

	err = message.Save(ctx, db.DB())
	if err != nil {
		return fmt.Errorf("Failed to save when creating message: %v", err)
	}

	theTrans := &models.Transaction{}
	lastTrans, err := theTrans.LastTransaction(ctx, db.DB(), req.UserID)
	if err != nil {
		if err != sql.ErrNoRows {
			return fmt.Errorf("Could not get user balance from db: %v", err)
		}
	}

	trans := &models.Transaction{
		UserID:   req.UserID,
		Amount:   req.Cost * -1,
		Currency: req.Currency,
		Balance:  lastTrans.Balance - req.Cost,
		Code:     strconv.FormatInt(message.ID, 10),
		Type:     "outgoing_message",
	}

	err = trans.Save(ctx, db.DB())
	if err != nil {
		return fmt.Errorf("Failed to save when creating transaction: %v", err)
	}

	for _, rec := range req.Recipients {
		recipient := &models.Recipient{
			MessageID:  message.ID,
			Phone:      rec.Number,
			Cost:       rec.Cost,
			Route:      req.Route,
			Currency:   req.Currency,
			Correlator: rec.MessageID,
		}

		err = recipient.Save(ctx, db.DB())
		if err != nil {
			return fmt.Errorf("Failed to save when creating recipient: %v", err)
		}

		newDlr := &models.Dlr{
			RecipientID: recipient.ID,
			Status:      rec.Status,
		}
		err = newDlr.Save(ctx, db.DB())
		if err != nil {
			revel.AppLog.Errorf("failed to save dlr: %v", err)
		}

		if len(req.StatusURL) > 0 {
			if rec.Status == "Success" {
				dlrReq := &entities.DLRRequest{
					ID:     rec.MessageID,
					Status: rec.Status,
					Source: req.Route,
					URL:    req.StatusURL,
				}

				err := h.processDlr(ctx, dlrReq)
				if err != nil {
					return fmt.Errorf("could not queue dlr task: %v", err)
				}
			}
		}
	}
	return nil
}

func (h *ProcessSMSJobHandler) processDlr(
	ctx context.Context,
	dlrReq *entities.DLRRequest,
) error {
	rand.Seed(time.Now().UnixNano())

	_, err := h.jobEnqueuer.EnqueueIn(
		ctx,
		sms_jobs.NewProcessDlrJob(dlrReq),
		time.Second*time.Duration(rand.Intn(600)),
	)
	if err != nil {
		return fmt.Errorf("could not queue dlr: %v", err)
	}

	return nil
}
