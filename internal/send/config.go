package send

import (
	"github.com/spf13/viper"
)

func Prepare() {
	viper.SetDefault("endpoint-uri", "http://localhost:8080")
	viper.SetDefault("logging.json", false)
	viper.SetDefault("logging.level", "warning")
}
