package rsq

import (
	"log"
	"os"
	"testing"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func TestSqsQueueAdapter(t *testing.T) {
	ran := false

	router := NewJobRouter()
	router.Handle("testing", func(job *Job) error {
		payload := string(job.Payload)
		if "testing" != payload {
			t.Errorf("Wrong payload got `%v`", payload)
		}

		ran = true
		return nil
	})

	options := SqsOptions{}
	options.QueueURL = os.Getenv("SQS_QUEUE")
	options.LongPollTimeout = 0
	options.MessagesPerWorker = 1

	queue := NewSqsAdapter(options)
	defer queue.Shutdown()

	err := queue.Push("testing", []byte("testing"))
	if err != nil {
		t.Fatal(err)
	}

	go queue.Work(router)

	time.Sleep(300 * time.Millisecond)

	if !ran {
		t.Fatal("The job should have been ran but wasn't")
	}
}
