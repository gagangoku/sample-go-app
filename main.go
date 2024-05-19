package main

import (
	"context"

	"github.com/rs/zerolog"
)

func main() {
	ctx := context.Background()
	logger := zerolog.Ctx(ctx)
	logger.Info().Msgf("Starting up: %s", VERSION)

	// Initialize app
	app := &App{}
	app.Init(ctx)

	// Start http server
	app.SetupHttpServer(ctx)
}
