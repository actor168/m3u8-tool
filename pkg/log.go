package pkg

import (
	log "github.com/sirupsen/logrus"
)

var LOGGER log.Logger

func Init() {
	LOGGER = log.New()
	LOGGER.Formatter = new(log.JSONFormatter)
	LOGGER.Formatter = new(log.TextFormatter)                     //default
	LOGGER.Formatter.(*log.TextFormatter).DisableColors = true    // remove colors
	LOGGER.Formatter.(*log.TextFormatter).DisableTimestamp = true // remove timestamp from test output
	LOGGER.Level = log.TraceLevel
	LOGGER.Out = os.Stdout
}
