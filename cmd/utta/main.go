package main

import (
	"os"

	"github.com/go-kit/kit/log"
)

func main() {
	logger := log.With(
		log.NewLogfmtLogger(os.Stderr),
		"ts",
		log.DefaultTimestampUTC,
	)

	app := NewApp(logger)
	if err := app.Command().Run(os.Args); err != nil {
		_ = app.Log("error", err)
		os.Exit(1)
	}
}
