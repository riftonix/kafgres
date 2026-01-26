package main

import (
	"kafgres/internal/app"
	"kafgres/internal/pkg/logger"
)

func main() {
	logger.Setup()
	app.Init()
}
