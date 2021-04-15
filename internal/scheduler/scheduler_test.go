package scheduler_test

import (
	"github.com/clambin/mediamon/internal/scheduler"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

type Task struct {
	lock sync.Mutex
	ran  bool
}

func (task *Task) Run() error {
	task.set()
	return nil
}

func (task *Task) set() {
	task.lock.Lock()
	defer task.lock.Unlock()
	task.ran = true
}

func (task *Task) get() bool {
	task.lock.Lock()
	defer task.lock.Unlock()
	return task.ran
}

func TestScheduler(t *testing.T) {
	s := scheduler.New()
	go s.Run()

	task := Task{}

	s.Schedule <- &scheduler.ScheduledTask{
		Task:     &task,
		Interval: 10 * time.Millisecond,
	}

	assert.Eventually(t, func() bool { return task.get() }, 100*time.Millisecond, 10*time.Millisecond)

	s.Stop <- struct{}{}
}
