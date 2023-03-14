package hystrix

import (
	"github.com/afex/hystrix-go/hystrix"
	"github.com/spf13/viper"
)

const CmdSample = "sample_command"

func init() {
	hystrix.ConfigureCommand(CmdSample, hystrix.CommandConfig{
		Timeout:                viper.GetInt("SAMPLE_TIMEOUT"),
		MaxConcurrentRequests:  viper.GetInt("SAMPLE_MAX_CONCURRENT_REQUESTS"),
		RequestVolumeThreshold: viper.GetInt("SAMPLE_REQUEST_VOLUME_THRESHOLD"),
		SleepWindow:            viper.GetInt("SAMPLE_SLEEP_WINDOW"),
		ErrorPercentThreshold:  viper.GetInt("SAMPLE_ERROR_PERCENT_THRESHOLD"),
	})
}
