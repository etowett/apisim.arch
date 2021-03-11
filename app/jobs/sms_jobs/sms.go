package sms_jobs

import (
	"apisim/app/entities"
	"apisim/app/jobs"
	"encoding/json"
)

const processSMSJobName = "process_sms_task"

type SendSMSJob struct {
	Request *entities.ProcessRequest `json:"request"`
}

func NewSendSMSJob(
	request *entities.ProcessRequest,
) *SendSMSJob {
	return &SendSMSJob{
		Request: request,
	}
}

func (h *SendSMSJob) JobName() string {
	return processSMSJobName
}

func (h *SendSMSJob) JobBody() (string, error) {
	b, err := json.Marshal(h)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (h *SendSMSJob) JobOptions() []jobs.PerformJobOption {
	return []jobs.PerformJobOption{
		jobs.WithMaxConcurrency(50),
		jobs.WithMaxFails(5),
		jobs.WithLowPriority(),
	}
}
