package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/spanner" // Required to initialize the spanner driver
	_ "github.com/golang-migrate/migrate/v4/source/file"      // Required to initialize the migration file
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const ReadWriteEveryonePermission = 0o666

var ErrMigrationFileAlreadyExist = errors.New("migration file already exist")

func GenerateMigrationFile(name string) error {
	for _, direction := range []string{"up", "down"} {
		basename := fmt.Sprintf("%s.%s.sql", name, direction)
		finalname := fmt.Sprintf("%s_%s", time.Now().Format("200601020304"), basename)

		matches, err := filepath.Glob("./db/migrations/*_" + basename)
		if err != nil {
			return err
		}

		if len(matches) > 0 {
			return ErrMigrationFileAlreadyExist
		}

		filename := filepath.Join("./db/migrations", finalname)

		_, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_EXCL, ReadWriteEveryonePermission)
		if err != nil {
			return err
		}

		absPath, _ := filepath.Abs(filename)
		log.Println(absPath)
	}

	log.Println("migrations file successfully created")
	return nil
}

// nolint
func Migrate(projectID, instanceID, databaseID string) error {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	if viper.GetString("SPANNER_EMULATOR_HOST") != "" {
		sugar.Infof("using spanner emulator at %s", viper.GetString("SPANNER_EMULATOR_HOST"))

		err := CreateInstanceAndDatabase(context.Background(), projectID, instanceID, databaseID)
		if err != nil {
			return err
		}
	}

	err := SetCredentials()
	if err != nil {
		return err
	}

	m, err := migrate.New(migrationPath(), buildDSN(projectID, instanceID, databaseID))
	if err != nil {
		sugar.Info("error while migrating db")
		return err
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	if errors.Is(err, migrate.ErrNoChange) {
		sugar.Info("nothing change. migration is up-to date")

		return nil
	}

	sugar.Info("migration completed")

	return nil
}

// nolint
func Rollback(projectID, instanceID, databaseName string) error {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	m, err := migrate.New(migrationPath(), buildDSN(projectID, instanceID, databaseName))
	if err != nil {
		sugar.Info("error while creating db rollback migration")

		return err
	}

	err = m.Steps(-1)
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		sugar.Info("error while running db rollback migration")

		return err
	}

	if errors.Is(err, migrate.ErrNoChange) {
		sugar.Info("nothing to rollback. migration is up-to date")

		return nil
	}

	sugar.Info("db rollback migration completed")
	return nil
}

func buildDSN(projectID, instanceID, databaseName string) string {
	return fmt.Sprintf(
		"spanner://projects/%s/instances/%s/databases/%s?%s",
		projectID,
		instanceID,
		databaseName,
		"x-clean-statements=true&x-migrations-table=schema_migrations",
	)
}

func migrationPath() string {
	if viper.GetString("APP_ENVIRONMENT") != "production" && viper.GetString("APP_ENVIRONMENT") != "integration" {
		_, currentFilePath, _, _ := runtime.Caller(0)

		dir := filepath.Dir(currentFilePath)
		return fmt.Sprintf("file://%s", filepath.Join(dir, "migrations"))
	}

	return "file://db/migrations"
}
