# rsq - Really Simple Queue

A simple queue system for go.  Provides an abstraction and implemntions
for a basic in memory queue and amazon's simple queue service.

```go get -u github.com/nerdyworm/rsq```

## Usage

``` go
router := rsq.NewJobRouter()
router.Handle("testing", func(job *Job) error {
  fmt.Printf("hello %s", job.Payload)
  return nil
})

queue := NewMemoryAdapter()
defer queue.Shutdown()

queue.Push("testing", []byte("testing"))
queue.Work(router)
```
