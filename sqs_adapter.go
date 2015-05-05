package rsq

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/service/sqs"
)

type SqsOptions struct {
	// AwsConfig is the config that you can pass into the sqs constructor
	AwsConfig *aws.Config

	// QueueURL is the amazon sqs url
	QueueURL string

	// MessagesPerWorker is the number of messages a worker can fetch at a time.
	// This adapter will also spawn a go routine per each message.
	MessagesPerWorker int64

	// LongPollTimeout how long to sit waiting for new messages from sqs
	LongPollTimeout int64
}

type SqsAdapter struct {
	service     *sqs.SQS
	ch          chan bool
	waitGroup   *sync.WaitGroup
	workerCount int

	queueURL          *string
	messagesPerWorker *int64
	longPollTimeout   *int64
}

func NewSqsAdapter(config SqsOptions) *SqsAdapter {
	return &SqsAdapter{
		ch:                make(chan bool),
		service:           sqs.New(config.AwsConfig),
		waitGroup:         &sync.WaitGroup{},
		queueURL:          aws.String(config.QueueURL),
		messagesPerWorker: aws.Long(config.MessagesPerWorker),
		longPollTimeout:   aws.Long(config.LongPollTimeout),
	}
}

func (s *SqsAdapter) Push(name string, payload []byte) error {
	message, err := s.newMessage(Job{name, payload})
	if err != nil {
		return err
	}

	_, err = s.service.SendMessage(message)
	if err != nil {
		return err
	}

	return nil
}

func (s *SqsAdapter) newMessage(job Job) (*sqs.SendMessageInput, error) {
	payload, err := json.Marshal(job)
	if err != nil {
		return nil, err
	}

	message := &sqs.SendMessageInput{
		QueueURL:    s.queueURL,
		MessageBody: aws.String(string(payload)),
	}

	return message, nil
}

func (s *SqsAdapter) Work(handler JobHandler) {
	s.waitGroup.Add(1)
	defer s.waitGroup.Done()

	s.workerCount += 1
	worker := s.workerCount
	log.Printf("starting worker %d\n", worker)

	for {
		select {
		case <-s.ch:
			log.Printf("stoping worker %d\n", worker)
			return
		default:
			messages, err := s.receiveMessages()
			if err != nil {
				log.Println(err)
				break
			}

			count := len(messages.Messages)
			errs := make(chan error, count)

			for _, message := range messages.Messages {
				go s.handleMessage(handler, message, errs)
			}

			for i := 0; i < count; i++ {
				e := <-errs
				if e != nil {
					log.Println(e)
				}
			}
		}
	}
}

func (s *SqsAdapter) receiveMessages() (*sqs.ReceiveMessageOutput, error) {
	input := &sqs.ReceiveMessageInput{
		QueueURL:            s.queueURL,
		MaxNumberOfMessages: s.messagesPerWorker,
		WaitTimeSeconds:     s.longPollTimeout,
	}

	return s.service.ReceiveMessage(input)
}

func (s *SqsAdapter) handleMessage(handler JobHandler, message *sqs.Message, status chan error) {
	job := Job{}

	err := json.Unmarshal([]byte(*message.Body), &job)
	if err != nil {
		status <- err
	}

	err = handler.Run(&job)
	if err == nil {
		err = s.deleteMessage(message)
	}

	status <- err
}

func (s *SqsAdapter) deleteMessage(message *sqs.Message) error {
	input := &sqs.DeleteMessageInput{
		QueueURL:      s.queueURL,
		ReceiptHandle: message.ReceiptHandle,
	}

	_, err := s.service.DeleteMessage(input)
	return err
}

func (s *SqsAdapter) Shutdown() error {
	log.Println("waiting for jobs to complete")
	close(s.ch)
	s.waitGroup.Wait()
	log.Println("shutdown succesful")
	return nil
}

var _ Queue = &SqsAdapter{}
