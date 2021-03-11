package jobs

type Job interface {
	JobName() string
	JobBody() (string, error) // Stored as JSON string
	JobOptions() []PerformJobOption
}
