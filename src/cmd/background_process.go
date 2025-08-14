package cmd

import (
	"github.com/spf13/cobra"
	"go-template/src/service"
)

var backgroundProcessCmd = &cobra.Command{
	Use:   "background-process",
	Short: "Use this command for run background process",

	RunE: func(cmd *cobra.Command, args []string) error {
		logger, err := getLogger()
		if err != nil {
			return err
		}

		service, err := service.NewService(logger)
		if err != nil {
			return err
		}

		return service.NewContext(nil).StartTimer()
	},
}

func init() {
	rootCmd.AddCommand(backgroundProcessCmd)
}
