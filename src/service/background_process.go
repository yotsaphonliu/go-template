package service

import (
	"github.com/go-co-op/gocron/v2"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (ctx *Context) StartTimer() error {
	ctx.Logger.Infof("User Management Backend E-Officer - Background Process Start !!")

	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		ctx.Logger.Warnf("Failed to load Asia/Bangkok timezone, using local timezone: %v", err)
		loc = time.Local
	}

	// สร้าง scheduler ใช้ timezone ที่ต้องการ
	s, err := gocron.NewScheduler(
		gocron.WithLocation(loc),
		gocron.WithLimitConcurrentJobs(2, gocron.LimitModeReschedule),
	)
	if err != nil {
		ctx.Logger.Errorf("Cannot create new scheduler: %v", err)
		return err
	}

	_, err = s.NewJob(
		gocron.DurationJob(time.Hour),
		gocron.NewTask(ctx.RemoveExpireApiKey),
	)
	if err != nil {
		ctx.Logger.Errorf("Cannot RemoveExpireApiKey job: %v", err)
		return err
	}

	s.Start()
	ctx.Logger.Infof("Background process scheduler started successfully")

	// Set up signal handling to gracefully shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for shutdown signal
	sig := <-sigChan
	ctx.Logger.Infof("Received signal %v, shutting down background process...", sig)

	// Shutdown the scheduler gracefully
	err = s.Shutdown()
	if err != nil {
		ctx.Logger.Errorf("Error during scheduler shutdown: %v", err)
		return err
	}

	ctx.Logger.Infof("Background process stopped gracefully")
	return nil

}
