package main

import (
	"log"
	"os"
	"time"

	"github.com/malamsyah/go-skele/cmd/server"
	"github.com/malamsyah/go-skele/db"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	cmd := &cobra.Command{
		Use:   "go-skele",
		Short: "A Skeleton project for Golang",
		Long:  "",
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:              "server",
			Short:            "HTTP Server Listener",
			Long:             "a http server listener that will be listening http request to this service.",
			TraverseChildren: true,
			RunE: func(cmd *cobra.Command, args []string) error {
				server.New(
					viper.GetString("APP_NAME"),
					viper.GetString("APP_PORT"),
					viper.GetDuration("GRACEFUL_SHUTDOWN_TIMEOUT_MS")*time.Millisecond).
					Start()

				return nil
			},
		},
		&cobra.Command{
			Use:   "generate-migration-file [file_name]",
			Short: "Generate db migration file",
			Long:  "Generate db migration file using date as file_name",
			RunE: func(_ *cobra.Command, args []string) error {
				return db.GenerateMigrationFile(args[0])
			},
		},
		&cobra.Command{
			Use:   "migrate",
			Short: "Running DB migration",
			Long:  "A DB migration command that will run DDL operation",
			RunE: func(_ *cobra.Command, _ []string) error {
				return db.Migrate(
					viper.GetString("SPANNER_PROJECT_ID"),
					viper.GetString("SPANNER_INSTANCE_ID"),
					viper.GetString("SPANNER_DATABASE_ID"),
				)
			},
		},
		&cobra.Command{
			Use:   "rollback",
			Short: "Running DB rollback",
			Long:  "A DB migration command that will rollback DDL operation",
			RunE: func(_ *cobra.Command, _ []string) error {
				return db.Rollback(
					viper.GetString("SPANNER_PROJECT_ID"),
					viper.GetString("SPANNER_INSTANCE_ID"),
					viper.GetString("SPANNER_DATABASE_ID"),
				)
			},
		})

	cobra.OnInitialize(bootstrap)
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}

func bootstrap() {
	loadConfig()
}

func loadConfig() {
	viper.AutomaticEnv()
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("can't load config from `.env`. environment variables will be used. err: %v", err)
	}

	// Set ENV to make Spanner connect to spanner emulator (used in development purpose only).
	// SPANNER_EMULATOR_HOST config on production env MUST BE empty
	_ = os.Setenv("SPANNER_EMULATOR_HOST", viper.GetString("SPANNER_EMULATOR_HOST"))
}
