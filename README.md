# rsq - Really Simple Queue

A simple queue system for go.  Provides an abstraction and implemntions
for a basic in memory queue and amazon's simple queue service.

```go get -u github.com/nerdyworm/rsq```

## Usage

``` go
package main

import (
	"fmt"

	"github.com/nerdyworm/rsq"
)

func main() {
	router := rsq.NewJobRouter()
	router.Handle("testing", func(job *rsq.Job) error {
		fmt.Printf("hello %s\n", job.Payload)
		return nil
	})

	queue := rsq.NewMemoryAdapter()
	defer queue.Shutdown()

	queue.Push("testing", []byte("world"))
	queue.Work(router)
}
```
