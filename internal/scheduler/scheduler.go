
package scheduler

import (
	"context"
	"fmt"
	"encoding/json"
	"velora/internal/workflow"

	"github.com/robfig/cron/v3"
)

// Scheduler handles scheduled tasks.
type Scheduler struct {
	cron   *cron.Cron
	engine *workflow.Engine
	jobs   map[cron.EntryID]string
}

// NewScheduler creates a new Scheduler.
func NewScheduler(engine *workflow.Engine) *Scheduler {
	return &Scheduler{
		cron:   cron.New(),
		engine: engine,
		jobs:   make(map[cron.EntryID]string),
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
		
		// Parse JSON to Workflow
		var wf workflow.Workflow
		if err := json.Unmarshal([]byte(jsonWorkflow), &wf); err != nil {
			fmt.Printf("Failed to parse workflow JSON: %v\n", err)
			return
		}
		
		// Run workflow
		ctx := context.Background()
		_, err := s.engine.Run(ctx, &wf, "")
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
