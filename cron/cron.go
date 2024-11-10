package cron

import (
	"github.com/robfig/cron"
	"verve/controller"
)

// Start cron with one scheduled job
func Start() {
	c := cron.New()
	err := c.AddFunc("0 * * * *", controller.CountLogger)
	if err != nil {
		return
	}
	c.Start()
}
