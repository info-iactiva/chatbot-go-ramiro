package handlers

import (
	"log"

	"github.com/robfig/cron/v3"
)

// StartScheduler inicia un scheduler genérico
func StartScheduler(schedule string, task func()) {
	c := cron.New()

	// Programar la tarea
	_, err := c.AddFunc(schedule, task)
	if err != nil {
		log.Fatalf("Error scheduling task: %v", err)
	}
	// task()

	c.Start()
	log.Printf("Scheduler started. Task scheduled for %s.\n", schedule)
}
