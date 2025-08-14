package migrations

import (
	"fmt"
	"sort"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go-template/src/core/log"
)

type Migration struct {
	Number uint `gorm:"primary_key"`
	Name   string

	Forwards func(db *gorm.DB) error `gorm:"-"`
}

var Migrations []*Migration

func Migrate(dryRun bool, number int, forceMigrate bool) error {
	dbHost := viper.GetString("Database.PostgreSQL.Host")
	dbPort := viper.GetString("Database.PostgreSQL.Port")
	dbUser := viper.GetString("Database.PostgreSQL.Username")
	sslMode := viper.GetString("Database.PostgreSQL.SSLMode")

	if dbUser == "" {
		dbUser = "postgres"
	}
	dbPassword := viper.GetString("Database.PostgreSQL.Password")
	if dbPassword == "" {
		dbPassword = "postgres"
	}
	dbName := viper.GetString("Database.PostgreSQL.DBName")

	if sslMode == "" {
		sslMode = "disable"
	}

	configLogger, err := log.InitConfig()
	if err != nil {
		return err
	}

	logger, err := log.NewLoggerWithModuleName(configLogger, "db/postgresql/migrations")
	if err != nil {
		return err
	}

	if dryRun {
		logger.Infof("=== DRY RUN ===")
	}

	// check for duplicate migration Number
	migrationIDs := make(map[uint]struct{})
	for _, migration := range Migrations {
		if _, ok := migrationIDs[migration.Number]; ok {
			err := fmt.Errorf("Duplicate migration Number found: %d", migration.Number)
			logger.Errorf("Unable to apply migrations, err: %+v", err)
			return err
		}

		migrationIDs[migration.Number] = struct{}{}
	}

	sort.Slice(Migrations, func(i, j int) bool {
		return Migrations[i].Number < Migrations[j].Number
	})

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost,
		dbPort,
		dbUser,
		dbPassword,
		dbName,
		sslMode,
	)
	db, err := gorm.Open("postgres", connStr)
	if err != nil {
		logger.Fatalf("%+v", err)
	}
	defer db.Close()

	// Force Migrate Zone
	if forceMigrate {
		logger.Infof("=== FORCE MIGRATE ===")
		if err := db.DropTableIfExists(&Migration{}).Error; err != nil {
			return errors.Wrap(err, "unable to drop migrations table")
		}
	}

	// Make sure Migration table is there
	logger.Debugf("ensuring migrations table is present")
	if err := db.AutoMigrate(&Migration{}).Error; err != nil {
		return errors.Wrap(err, "unable to automatically migrate migrations table")
	}

	var latest Migration
	if err := db.Order("number desc").First(&latest).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return errors.Wrap(err, "unable to find latest migration")
	}

	noMigrationsApplied := latest.Number == 0

	if noMigrationsApplied && len(Migrations) == 0 {
		logger.Infof("no migrations to apply")
		return nil
	}

	if latest.Number >= Migrations[len(Migrations)-1].Number {
		logger.Infof("no migrations to apply")
		return nil
	}

	if number == -1 {
		number = int(Migrations[len(Migrations)-1].Number)
	}

	if uint(number) <= latest.Number && latest.Number > 0 {
		logger.Infof("no migrations to apply, specified number is less than or equal to latest migration; backwards migrations are not supported")
		return nil
	}

	for _, migration := range Migrations {
		if migration.Number > uint(number) {
			break
		}

		if migration.Number <= latest.Number {
			continue
		}

		if latest.Number > 0 {
			logger.Infof("continuing migration starting from %d", migration.Number)
		}

		logger := logger.WithFields(log.Fields{
			"migration_number": migration.Number,
		})
		logger.Infof("applying migration %q", migration.Name)

		if dryRun {
			continue
		}

		tx := db.Begin()

		if err := migration.Forwards(tx); err != nil {
			logger.Errorf("unable to apply migration, rolling back. err: %+v", err)
			if err := tx.Rollback().Error; err != nil {
				logger.Errorf("unable to rollback... err: %+v", err)
			}
			break
		}

		if err := tx.Commit().Error; err != nil {
			logger.Errorf("unable to commit transaction... err: %+v", err)
			break
		}

		// Create migration record
		if err := db.Create(migration).Error; err != nil {
			logger.Errorf("unable to create migration record. err: %+v", err)
			break
		}
	}

	return nil
}
