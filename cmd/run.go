package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"signature.service/pkg/logger"
	"signature.service/pkg/services"
	storage_pkg "signature.service/pkg/storage"
	"signature.service/pkg/workspace"
)

func init() {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Runs the API server",
		Run: func(cmd *cobra.Command, args []string) {
			logger.NewMainLogger(
				logger.WithMinLevel(os.Getenv(logger.LOGGER_LEVEL_ENV)),
			)

			g, ctx := errgroup.WithContext(ctx)

			var storage storage_pkg.Storage
			switch os.Getenv("STORAGE") {
			case "memory":
				storage = storage_pkg.NewMemoryStorage()
			default:
				logger.Main.WithFields(map[string]any{"storage": os.Getenv("STORAGE")}).Fatal("unknown storage type")
			}

			// workspace factory
			workspaceFactory := workspace.NewWorkspaceFactory(storage)

			// api server
			serverCfg := &services.ServerConfig{
				Port: fmt.Sprintf(":%s", os.Getenv("API_SERVER_PORT")),
			}
			server := services.NewServer(ctx, serverCfg, workspaceFactory)

			g.Go(func() error { return server.Run(ctx, serverCfg.Port) })
			if err := g.Wait(); err != nil {
				logger.Main.WithError(err).Panic("failed to run")
			}
		},
	}

	rootCmd.AddCommand(runCmd)
}
