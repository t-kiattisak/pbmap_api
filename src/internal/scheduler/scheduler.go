package scheduler

import (
	"context"
	"fmt"
	"pbmap_api/src/internal/usecase"

	"github.com/robfig/cron/v3"
)

type Scheduler interface {
	Start()
	Stop()
}

type appScheduler struct {
	cron            *cron.Cron
	dataSyncService usecase.DataSyncService
}

func NewScheduler(dataSyncService usecase.DataSyncService) Scheduler {
	c := cron.New()

	s := &appScheduler{
		cron:            c,
		dataSyncService: dataSyncService,
	}

	s.registerJobs()
	return s
}

func (s *appScheduler) Start() {
	fmt.Println("[Scheduler] Starting cron scheduler...")
	s.cron.Start()
}

func (s *appScheduler) Stop() {
	fmt.Println("[Scheduler] Stopping cron scheduler...")
	s.cron.Stop()
}

func (s *appScheduler) registerJobs() {
	// Cron Spec: "@every 10s"
	_, err := s.cron.AddFunc("@every 10m", func() {
		ctx := context.Background()
		if err := s.dataSyncService.SyncAndNotify(ctx); err != nil {
			fmt.Printf("[Scheduler] Job failed: %v\n", err)
		}
	})

	if err != nil {
		fmt.Printf("[Scheduler] Failed to add job: %v\n", err)
	} else {
		fmt.Println("[Scheduler] Registered 'DataSync' job: @every 10s")
	}
}
