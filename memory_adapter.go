package rsq

import "log"

// NewMemoryAdapter returns an instance of a MemoryAdapter
func NewMemoryAdapter() *MemoryAdapter {
	return &MemoryAdapter{make([]Job, 0)}
}

// MemoryAdapter is a simple in memory queue system
type MemoryAdapter struct {
	jobs []Job
}

// Push a job on to the queue
func (q *MemoryAdapter) Push(name string, payload []byte) error {
	q.jobs = append(q.jobs, Job{name, payload})
	return nil
}

// Work off all the jobs in the queue
func (q *MemoryAdapter) Work(handler JobHandler) {
	for i, job := range q.jobs {
		err := handler.Run(&job)
		if err != nil {
			log.Printf("error running job %v\n", err)
		}

		if err == nil {
			// Remove job from queue
			q.jobs = append(q.jobs[:i], q.jobs[i+1:]...)
		}
	}
}

// Shutdown the queue and clean up all resources
func (q *MemoryAdapter) Shutdown() error {
	return nil
}

var _ Queue = NewMemoryAdapter()
