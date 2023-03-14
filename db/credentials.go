package db

import (
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// SetCredentials used for migration only
// it sets GOOGLE_APPLICATION_CREDENTIALS into env
// see https://cloud.google.com/docs/authentication/getting-started#setting_the_environment_variable
// this is required as github.com/golang-migrate/migrate for spanner assume that this env is being set
// nolint
func SetCredentials() error {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	if viper.GetString("SPANNER_EMULATOR_HOST") != "" {
		return nil
	}

	logger.Info("Saving credentials to a json and set to env")

	credentials := viper.GetString("SPANNER_CREDENTIALS")

	err := os.WriteFile("credentials.json", []byte(credentials), ReadWriteEveryonePermission)
	if err != nil {
		return err
	}

	err = os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "credentials.json")
	if err != nil {
		return err
	}

	return nil
}
