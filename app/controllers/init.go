package controllers

import (
	"apisim/app/db"
	"apisim/app/jobs"
	"apisim/app/jobs/job_handlers"
	"apisim/app/providers"
	"apisim/app/work"
	"context"
	"time"

	"github.com/revel/revel"
)

var (
	jobEnqueuer  *work.AppJobEnqueuer
	redisManager *db.AppRedis
	csvCreator   *providers.SimpleCSVCreator
)

func init() {
	revel.OnAppStart(initApp)
	revel.InterceptMethod((*App).AddUser, revel.BEFORE)
}

func initApp() {
	redisManager = db.NewRedisProvider(&db.RedisConfig{
		IdleTimeout: 2 * time.Minute,
		MaxActive:   1000,
		MaxIdle:     100,
	})

	redisPool := redisManager.RedisPool()
	jobEnqueuer = work.NewJobEnqueuer(redisPool)
	workerPool := work.NewWorkerPool(redisPool, uint(200))
	jobHandlers := setupJobHandlers(jobEnqueuer)
	workerPool.RegisterJobs(jobHandlers...)
	workerPool.Start(context.Background())
}

func setupJobHandlers(
	jobEnqueuer work.JobEnqueuer,
) []jobs.JobHandler {
	processSMSJobHandler := job_handlers.NewProcessSMSJobHandler(jobEnqueuer)
	processDlrJobHandler := job_handlers.NewProcessDlrJobHandler(jobEnqueuer)
	return []jobs.JobHandler{
		processDlrJobHandler,
		processSMSJobHandler,
	}
}
