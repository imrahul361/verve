package main

import (
	"verve/cron"
	"verve/router"
)

func main() {
	finish := make(chan bool)
	router.Init()
	cron.Start()
	<-finish
}
