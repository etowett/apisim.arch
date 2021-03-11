package work

import (
	"apisim/app/jobs"
	"context"
	"fmt"
	"time"

	gowork "github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"github.com/revel/revel"
)

type WorkerPool interface {
	RegisterJobs(jobHandlers ...jobs.JobHandler)
	Start(ctx context.Context)
	Stop()
}

type AppWorkerPool struct {
	ctx  context.Context
	pool *gowork.WorkerPool
}

type workerPoolContext struct{}

func log(
	job *gowork.Job,
	next gowork.NextMiddlewareFunc,
) error {
	revel.AppLog.Infof("Starting job id: %s name: %s", job.ID, job.Name)
	return next()
}

func NewWorkerPool(
	redisPool *redis.Pool,
	concurrency uint,
) *AppWorkerPool {
	pool := gowork.NewWorkerPool(workerPoolContext{}, concurrency, workerNamespace, redisPool)
	pool.Middleware(log)

	workerPool := &AppWorkerPool{
		pool: pool,
	}

	return workerPool
}

func (wp *AppWorkerPool) RegisterJobs(jobHandlers ...jobs.JobHandler) {

	for _, jobHandler := range jobHandlers {
		// Set up default options
		jobOptions := gowork.JobOptions{
			Priority: 1,
			MaxFails: 1,
		}

		job := jobHandler.Job()
		opts := job.JobOptions()
		if opts != nil {
			for _, opt := range opts {
				opt(jobOptions)
			}
		}

		wrappedJobHandler := wp.wrapJobHandler(jobHandler)

		wp.pool.JobWithOptions(job.JobName(), jobOptions, wrappedJobHandler)
	}
}

func (wp *AppWorkerPool) Start(
	ctx context.Context,
) {
	wp.ctx = ctx
	wp.pool.Start()
}

func (wp *AppWorkerPool) Stop() {
	wp.pool.Stop()
}

func (wp *AppWorkerPool) wrapJobHandler(jobHandler jobs.JobHandler) func(job *gowork.Job) error {
	return func(job *gowork.Job) error {
		startTime := time.Now()

		// var requestID string
		// rawContext, ok := job.Args[contextKey]
		// if ok {
		// 	context, ok := rawContext.(map[string]interface{})
		// 	if ok {
		// 		rawRequestID := context[requestIDKey]
		// 		requestID, ok = rawRequestID.(string)
		// 		if !ok {
		// 			requestID = ""
		// 		}
		// 	}
		// }

		// jobContext := ctxhelper.WithRequestId(wp.ctx, requestID)
		jobContext := context.WithValue(wp.ctx, "request-id", 1)

		rawBody := job.Args[bodyKey]
		if rawBody == nil {
			return nil
		}

		body, ok := rawBody.(string)
		if !ok {
			duration := time.Now().Sub(startTime)
			err := fmt.Errorf("failed to cast to string: %v", rawBody)
			revel.AppLog.Errorf("[JobHandler] Job failed: %v, job_id: %v, job_name: %v, duration: %v, args: %v", err, job.ID, job.Name, duration, job.Args)

			// Since this body cannot be handled properly, return nil to prevent retries
			return nil
		}

		if len(body) < 0 {
			return nil
		}

		err := jobHandler.PerformJob(jobContext, body)
		duration := time.Now().Sub(startTime)
		if err != nil {
			revel.AppLog.Errorf("[JobHandler] Job failed: %v, job_id: %v, job_name: %v, duration: %v, args: %v", err, job.ID, job.Name, duration, job.Args)
		} else {
			revel.AppLog.Infof("[JobHandler] Job completed, job_id: %v, job_name: %v, duration: %v, args: %v", job.ID, job.Name, duration, job.Args)
		}

		return err
	}
}
