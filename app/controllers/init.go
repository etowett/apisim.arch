package controllers

import (
	"apisim/app/db"
	"time"

	"github.com/revel/revel"
)

var (
	redisManager *db.AppRedis
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

	// redisPool := redisManager.RedisPool()

	// jobEnqueuer = work.NewJobEnqueuer(redisPool)

	// workerPool := work.NewWorkerPool(redisPool, uint(200))

	// jobHandlers := setupJobHandlers(
	// 	africasTalkingSender,
	// 	jobEnqueuer,
	// )

	// workerPool.RegisterJobs(jobHandlers...)

	// workerPool.Start(context.Background())
}
