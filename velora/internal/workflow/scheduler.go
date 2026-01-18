
package workflow

import (
	"sync"
	"time"
)

type ScheduledTask struct {
	Pipeline  *Pipeline
	Interval  time.Duration
	stop      chan struct{}
	isRunning bool
	mu        sync.Mutex
}

func (t *ScheduledTask) Start(r *Runner) {
	t.mu.Lock()
	if t.isRunning {
		t.mu.Unlock()
		return
	}
	t.isRunning = true
	t.mu.Unlock()

	t.stop = make(chan struct{})
	ticker := time.NewTicker(t.Interval)

	go func() {
		for {
			select {
			case <-ticker.C:
				// Execute the pipeline
				r.Run(t.Pipeline)
			case <-t.stop:
				ticker.Stop()
				return
			}
		}
	}()
}

func (t *ScheduledTask) Stop() {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.isRunning {
		close(t.stop)
		t.isRunning = false
	}
}


type Scheduler struct {
	tasks map[string]*ScheduledTask
	mu    sync.Mutex
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		tasks: make(map[string]*ScheduledTask),
	}
}

func (s *Scheduler) Schedule(name string, pipeline *Pipeline, interval time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task := &ScheduledTask{
		Pipeline: pipeline,
		Interval: interval,
	}
	s.tasks[name] = task
}

func (s *Scheduler) Start(name string, r *Runner) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if task, ok := s.tasks[name]; ok {
		task.Start(r)
	}
}

func (s *Scheduler) Stop(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if task, ok := s.tasks[name]; ok {
		task.Stop()
	}
}
