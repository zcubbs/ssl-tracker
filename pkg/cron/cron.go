package cron

import (
	"context"
	"github.com/charmbracelet/log"
	"github.com/robfig/cron/v3"
)

func StartCronJob(cronPattern string, task func(ctx context.Context)) {
	ctx := context.Background()

	if cronPattern == "" {
		log.Info("no cron pattern provided, not starting cron job")
		return
	}

	if cronPattern == "-" {
		log.Info("running cron job once")
		task(ctx)
		log.Info("cron job finished")
		return
	}

	c := cron.New(cron.WithSeconds()) // cron with second-level precision
	_, err := c.AddFunc(cronPattern, func() {
		task(ctx)
	})
	if err != nil {
		log.Fatalf("cannot create cron job: %v", err)
	}

	log.Info("starting cron job")
	c.Start()
}
