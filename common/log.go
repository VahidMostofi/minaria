package common

import (
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var logger *log.Logger

func GetLogger() *log.Logger {
	if logger == nil {
		logger = log.New()
		if viper.GetBool(DISABLE_LOGGING) {
			logger.SetOutput(ioutil.Discard)
		}
	}

	return logger
}
