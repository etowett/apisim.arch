package sms_jobs

import (
	"apisim/app/entities"
	"apisim/app/jobs"
	"encoding/json"
)

const processDlrJobName = "process_dlr_task"

type ProcessDlrJob struct {
	Request *entities.DLRRequest `json:"dlr_request"`
}

func NewProcessDlrJob(
	request *entities.DLRRequest,
) *ProcessDlrJob {
	return &ProcessDlrJob{
		Request: request,
	}
}

func (h *ProcessDlrJob) JobName() string {
	return processDlrJobName
}

func (h *ProcessDlrJob) JobBody() (string, error) {
	b, err := json.Marshal(h)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (h *ProcessDlrJob) JobOptions() []jobs.PerformJobOption {
	return []jobs.PerformJobOption{
		jobs.WithMaxConcurrency(50),
		jobs.WithMaxFails(5),
		jobs.WithLowPriority(),
	}
}
