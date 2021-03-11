package jobs

import (
	"github.com/gocraft/work"
)

const (
	defaultMaxConcurrency = 10
	defaultMaxFails       = 5
	defaultPriority       = 3
	lowPriority           = 1
	highPriority          = 10
)

type PerformJobOption func(work.JobOptions)

func WithMaxConcurrency(maxConcurrency uint) PerformJobOption {
	return func(jobOptions work.JobOptions) {
		jobOptions.MaxConcurrency = maxConcurrency
	}
}

func WithMaxFails(maxFails uint) PerformJobOption {
	return func(jobOptions work.JobOptions) {
		jobOptions.MaxFails = maxFails
	}
}

func WithLowPriority() PerformJobOption {
	return WithPriority(lowPriority)
}

func WithHighPriority() PerformJobOption {
	return WithPriority(highPriority)
}

func WithPriority(priority uint) PerformJobOption {
	return func(jobOptions work.JobOptions) {
		jobOptions.Priority = priority
	}
}

func WithSkipDeadQueue(skip bool) PerformJobOption {
	return func(jobOptions work.JobOptions) {
		jobOptions.SkipDead = skip
	}
}
