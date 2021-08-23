package main

import (
	"github.com/aerostatka/banking/app"
	"github.com/aerostatka/banking-lib/logger"
)

func main() {
	logger.Info("Starting application")
	app.Start()
}
