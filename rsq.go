package rsq

import "errors"

var (
	// ErrNoHandlerFound is returned from the default NotFoundHandler
	// when we can not match a Job to it's handler
	ErrNoHandlerFound = errors.New("not found")
)

// NewJobRouter returns a new router instance.
func NewJobRouter() *JobRouter {
	return &JobRouter{
		NotFoundHandler: defaultNotFound,
		namedJobs:       make(map[string]*JobRoute),
	}
}

// JobRouter registers routes to be matched and dispathces a handler.
type JobRouter struct {
	// NotFoundHandler will be called when a handler can not be found for a job
	NotFoundHandler JobHandlerFunc

	// Jobs by name
	namedJobs map[string]*JobRoute
}

// JobRoute stores information to match incoming jobs
type JobRoute struct {
	name string
	fn   JobHandlerFunc
}

// Job is the datastructure passed in to JobHandlers
type Job struct {
	Name    string
	Payload []byte
}

// JobHandlerFunc handles a job
//
// It should return an error if the job could not
// be ran sucessfully.
//
// It should retur nil if the job was succesful and
// can be removed from the Queue
type JobHandlerFunc func(job *Job) error

// JobHandler handler jobs
type JobHandler interface {
	Run(*Job) error
}

var _ JobHandler = &JobRouter{}

// Handle adds a new JobHandlerFunc to the router
func (r *JobRouter) Handle(name string, fn JobHandlerFunc) {
	route := &JobRoute{name, fn}
	r.namedJobs[name] = route
}

// Run a job
func (r *JobRouter) Run(job *Job) error {
	if handler, ok := r.namedJobs[job.Name]; ok {
		return handler.fn(job)
	}

	return r.NotFoundHandler(job)
}

func defaultNotFound(job *Job) error {
	return ErrNoHandlerFound
}
