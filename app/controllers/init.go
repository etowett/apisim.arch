package controllers

import (
	"apisim/app/db"
	"apisim/app/jobs"
	"apisim/app/jobs/job_handlers"
	"apisim/app/providers"
	"apisim/app/work"
	"context"
	"fmt"
	"time"

	"github.com/revel/revel"
	"golang.org/x/text/message"
)

var (
	jobEnqueuer  *work.AppJobEnqueuer
	redisManager *db.AppRedis
	csvCreator   *providers.SimpleCSVCreator
)

func init() {
	revel.OnAppStart(initApp)
	revel.InterceptMethod((*App).AddUser, revel.BEFORE)
	revel.InterceptMethod(Outbox.checkUser, revel.BEFORE)
	revel.InterceptMethod(Settings.checkUser, revel.BEFORE)

	revel.TemplateFuncs["formatDate"] = func(theTime time.Time) string {
		timeLocation, err := time.LoadLocation("Africa/Nairobi")
		if err != nil {
			revel.AppLog.Errorf("failed to load Nairobi timezone: %+v", err)
			return theTime.Format("Jan _2 2006 3:04PM")
		}

		return theTime.In(timeLocation).Format("Jan _2 2006 3:04PM")
	}

	revel.TemplateFuncs["formatMoney"] = func(currency string, amount float64) string {
		p := message.NewPrinter(message.MatchLanguage("en"))
		return fmt.Sprintf(p.Sprintf("%v %.2f", currency, amount))
	}
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
