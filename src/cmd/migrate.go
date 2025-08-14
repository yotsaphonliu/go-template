package cmd

import (
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/cobra"
	"go-template/src/core/db"
)

var migrateDBCmd = &cobra.Command{
	Use: "migrate-db",
	RunE: func(cmd *cobra.Command, args []string) error {
		number, _ := cmd.Flags().GetInt("number")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		forceMigrate, _ := cmd.Flags().GetBool("force-migrate")

		dbConfig, err := db.InitConfig()
		if err != nil {
			return err
		}

		db.Migrate(dbConfig, dryRun, number, forceMigrate)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(migrateDBCmd)

	migrateDBCmd.Flags().Int("number", -1, "the migration to run forwards until; if not set, will run all migrations")
	migrateDBCmd.Flags().Bool("dry-run", false, "print out migrations to be applied without running them")
	migrateDBCmd.Flags().Bool("force-migrate", false, "drop all the tables before migrate the database")
}
