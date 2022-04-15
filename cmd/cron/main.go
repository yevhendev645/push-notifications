package main

import (
	"log"
	"notifications/taskmgr"
	"notifications/tasks"
	"os"
	"os/signal"

	"github.com/robfig/cron/v3"
)

func main() {
	c := cron.New()

	// Run cron every 1 minute
	c.AddFunc("@every 1m", func() {
		taskmgr.EnqueueTask(&tasks.CheckBalanceChange{}, taskmgr.TaskOptions{})
		taskmgr.EnqueueTask(&tasks.CheckCryptoChange{}, taskmgr.TaskOptions{})
	})

	c.Start()
	log.Println("=====cron system started======")

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	<-sigint
}
