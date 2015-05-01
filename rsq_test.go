package rsq

import (
	"errors"
	"testing"
)

func TestJobRouter(t *testing.T) {
	router := NewJobRouter()

	err := router.Run(&Job{})
	if err != ErrNoHandlerFound {
		t.Fatal("An empty router should return a ErrNoHandlerFound error")
	}

	e := errors.New("my custom error")
	router.NotFoundHandler = func(j *Job) error {
		return e
	}

	err = router.Run(&Job{})
	if e != err {
		t.Fatal("Did not return the error from the NotFoundHandler")
	}
}
