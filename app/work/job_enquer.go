package work

import (
	"apisim/app/jobs"
	"context"
	"time"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"github.com/revel/revel"
)

const (
	workerNamespace = "apisim"
	contextKey      = "_context"
	requestIDKey    = "request_id"
	bodyKey         = "body"
)

type JobEnqueuer interface {
	Enqueue(ctx context.Context, job jobs.Job) (string, error)
	EnqueueIn(ctx context.Context, job jobs.Job, duration time.Duration) (string, error)
}

type AppJobEnqueuer struct {
	enqueuer *work.Enqueuer
}

func NewJobEnqueuer(pool *redis.Pool) *AppJobEnqueuer {
	return &AppJobEnqueuer{
		enqueuer: work.NewEnqueuer(workerNamespace, pool),
	}
}

func (e *AppJobEnqueuer) Enqueue(ctx context.Context, job jobs.Job) (string, error) {
	return e.enqueue(ctx, job, 0)
}

func (e *AppJobEnqueuer) EnqueueIn(ctx context.Context, job jobs.Job, duration time.Duration) (string, error) {
	return e.enqueue(ctx, job, duration)
}

func (e *AppJobEnqueuer) enqueue(
	ctx context.Context,
	job jobs.Job,
	duration time.Duration,
) (string, error) {

	b, err := job.JobBody()
	if err != nil {
		return "", err
	}

	args := make(map[string]interface{})
	contextArgs := e.contextArgs(ctx)
	args[contextKey] = contextArgs
	args[bodyKey] = b

	var internalJob *work.Job
	if duration > 0 {
		var scheduledJob *work.ScheduledJob
		scheduledJob, err = e.enqueuer.EnqueueIn(job.JobName(), int64(duration.Seconds()), args)
		if err != nil {
			internalJob = scheduledJob.Job
		}
	} else {
		internalJob, err = e.enqueuer.Enqueue(job.JobName(), args)
	}

	if err != nil {
		revel.AppLog.Errorf("[JobEnqueuer] Failed to enqueue job, job_name: %v, duration: %v, body: %v, context: %v, err: %v", job.JobName(), duration, string(b), contextArgs, err)
		return "", err
	}

	revel.AppLog.Infof("[JobEnqueuer] Successfully enqueued job, job_id: %v, job_name: %v, body: %v, context: %v", internalJob.ID, job.JobName(), string(b), contextArgs)

	return internalJob.ID, nil
}

func (e *AppJobEnqueuer) contextArgs(ctx context.Context) map[string]interface{} {
	args := make(map[string]interface{})
	args[requestIDKey] = "1"
	// args[requestIDKey] = ctxhelper.RequestId(ctx)
	return args
}
