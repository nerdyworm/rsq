package rsq

// Queue is an implemention of a queue
type Queue interface {
	// Push adds a job to this queue
	Push(name string, payload []byte) error

	// Work should take jobs from it's queue and pass them into the JobRouter
	// It is up to the implemention to decide if this should loop forever or
	// just work off the jobs that are currently in the queue
	Work(JobHandler)

	// Shutdown should free any resources and clear connections.  It should
	// also wait until all jobs are done processing.
	Shutdown() error
}
