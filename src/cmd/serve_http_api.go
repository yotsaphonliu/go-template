package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"go-template/src/core/handlers/routes"
	"go-template/src/otel"
	"go-template/src/service"
)

var serveAPICmd = &cobra.Command{
	Use:   "serve-http-api",
	Short: "Start HTTP API server",
	RunE: func(cmd *cobra.Command, args []string) error {

		logger, err := getLogger()
		if err != nil {
			return err
		}

		service, err := service.NewService(logger)
		if err != nil {
			return err
		}

		config, err := routes.InitConfig()
		if err != nil {
			return err
		}

		ctx := context.Background()
		tp, err := otel.Init(ctx)
		if err != nil {
			return err
		}
		defer func() {
			if err := tp.Shutdown(ctx); err != nil {
				logger.Errorf("Failed to shutdown tracer provider: %v", err)
			}
		}()

		routes.NewRouter(config, logger, service)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(serveAPICmd)
}
