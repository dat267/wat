package main

import (
	"context"
	"time"
)

type Job struct {
	Name     string
	Interval time.Duration
	Action   func(a *App) (string, error)
}

var schedulerCtx context.Context
var schedulerCancel context.CancelFunc

func startScheduler(a *App) {
	schedulerCtx, schedulerCancel = context.WithCancel(context.Background())
	jobs := []Job{
		{
			Name:     "Disk Check",
			Interval: 5 * time.Minute,
			Action: func(app *App) (string, error) {
				return app.ExecuteScript("disk")
			},
		},
		{
			Name:     "Network Check",
			Interval: 3 * time.Minute,
			Action: func(app *App) (string, error) {
				return app.ExecuteScript("network")
			},
		},
	}
	for _, j := range jobs {
		go runJobLoop(a, j)
	}
}

func stopScheduler() {
	if schedulerCancel != nil {
		schedulerCancel()
	}
}

func runJobLoop(a *App, j Job) {
	ticker := time.NewTicker(j.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-schedulerCtx.Done():
			return
		case <-ticker.C:
			runJob(a, j)
		}
	}
}

func runJob(a *App, j Job) {
	timestamp := time.Now().Format(time.RFC3339)
	logInfo("Running scheduled task: "+j.Name, nil)
	out, err := j.Action(a)
	status := "success"
	detail := out
	if err != nil {
		status = "error"
		detail = err.Error()
		logError("Scheduled task failed: "+j.Name, "error", err)
	}
	a.StorePut("scheduler_history", j.Name+"_"+timestamp, status+"|"+detail)
}
