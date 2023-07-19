package tls

import (
	"context"
	"github.com/charmbracelet/log"
	"github.com/robfig/cron/v3"
)

func StartCheckCertificateValidityCronJob(cronPattern string, task func(ctx context.Context)) {
	ctx := context.Background()

	if cronPattern == "" {
		log.Info("No cron pattern provided, not starting cron job")
		return
	}

	if cronPattern == "-" {
		log.Info("Checking certificate validity once")
		task(ctx)
		log.Info("Done checking certificate validity")
		return
	}

	c := cron.New(cron.WithSeconds()) // cron with second-level precision
	_, err := c.AddFunc(cronPattern, func() {
		task(ctx)
	})
	if err != nil {
		log.Fatalf("Cannot create cron job: %v", err)
	}

	log.Info("Starting cron job")
	c.Start()
}
