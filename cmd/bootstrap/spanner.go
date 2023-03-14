package bootstrap

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/spanner"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
)

func BuildSpannerClient() *spanner.Client {
	dsn := fmt.Sprintf(
		"projects/%s/instances/%s/databases/%s",
		viper.GetString("SPANNER_PROJECT_ID"),
		viper.GetString("SPANNER_INSTANCE_ID"),
		viper.GetString("SPANNER_DATABASE_ID"),
	)

	// if SPANNER_EMULATOR_HOST is not set, we're running in integration or production
	if viper.GetString("SPANNER_EMULATOR_HOST") == "" {
		credentials := viper.GetString("SPANNER_CREDENTIALS")

		db, err := spanner.NewClient(context.Background(), dsn, option.WithCredentialsJSON([]byte(credentials)))
		if err != nil {
			log.Fatalf("failed building repository. spanner.NewClient %v", err)
		}

		return db
	}

	db, err := spanner.NewClient(context.Background(), dsn)
	if err != nil {
		log.Fatalf("failed building repository. spanner.NewClient %v", err)
	}

	return db
}
