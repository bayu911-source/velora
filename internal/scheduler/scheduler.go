
package scheduler

import (
	"fmt"
	"velora/internal/workflow"

	"github.com/robfig/cron/v3"
)

// Scheduler handles scheduled tasks.
type Scheduler struct {
	cron           *cron.Cron
	pipelineRunner *workflow.PipelineRunner
	jobs           map[cron.EntryID]string
}

// NewScheduler creates a new Scheduler.
func NewScheduler(pipelineRunner *workflow.PipelineRunner) *Scheduler {
	return &Scheduler{
		cron:           cron.New(),
		pipelineRunner: pipelineRunner,
		jobs:           make(map[cron.EntryID]string),
	}
}

// Start starts the cron scheduler.
func (s *Scheduler) Start() {
	s.cron.Start()
}

// Stop stops the cron scheduler.
func (s *Scheduler) Stop() {
	s.cron.Stop()
}

// AddWorkflow schedules a workflow to run at a specific cron interval.
func (s *Scheduler) AddWorkflow(spec, jsonWorkflow string) (cron.EntryID, error) {
	job := func() {
		fmt.Printf("Running workflow %s...\n", jsonWorkflow)
		_, err := s.pipelineRunner.Run(jsonWorkflow)
		if err != nil {
			fmt.Printf("Workflow failed: %v\n", err)
		}
	}

	id, err := s.cron.AddFunc(spec, job)
	if err != nil {
		return 0, err
	}

	s.jobs[id] = jsonWorkflow
	return id, nil
}

// ListJobs returns a map of scheduled jobs.
func (s *Scheduler) ListJobs() map[cron.EntryID]string {
	return s.jobs
}
