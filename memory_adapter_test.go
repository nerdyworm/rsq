package rsq

import (
	"errors"
	"testing"
)

func TestMemoryAdapter(t *testing.T) {
	ran := false

	router := NewJobRouter()
	router.Handle("testing", func(job *Job) error {
		t.Log(job)
		ran = true
		return nil
	})

	queue := NewMemoryAdapter()
	defer queue.Shutdown()

	queue.Push("testing", []byte("testing"))
	queue.Work(router)

	if !ran {
		t.Fatal("The job should have been ran but wasn't")
	}
}

func TestMemoryAdapterRemovesRanJobs(t *testing.T) {
	count := 0

	router := NewJobRouter()
	router.Handle("testing", func(job *Job) error {
		count++
		return nil
	})

	queue := NewMemoryAdapter()
	defer queue.Shutdown()
	queue.Push("testing", []byte("testing"))
	queue.Work(router)
	queue.Work(router)

	if count > 1 {
		t.Fatalf("The job should have only been ran once it was ran %d", count)
	}
}

func TestMemoryAdapterDoesNotRemoveFailedJob(t *testing.T) {
	count := 0

	router := NewJobRouter()
	router.Handle("testing", func(job *Job) error {
		count++
		return errors.New("don't remove me")
	})

	queue := NewMemoryAdapter()
	defer queue.Shutdown()

	queue.Push("testing", []byte("testing"))
	queue.Work(router)
	queue.Work(router)

	if count != 2 {
		t.Fatalf("The job should have been ran twice it was ran %d", count)
	}
}
