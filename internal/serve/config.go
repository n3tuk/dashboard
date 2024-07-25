package serve

import (
	"github.com/spf13/viper"
)

const (
	host = "localhost"
	port = 8080
)

func Prepare() {
	viper.SetDefault("web.bind.hostname", host)
	viper.SetDefault("web.bind.port", port)
	viper.SetDefault("logging.json", true)
	viper.SetDefault("logging.level", "info")
}
