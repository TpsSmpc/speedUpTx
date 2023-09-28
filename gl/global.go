package gl

import (
	"os"

	"github.com/triplefi/go-logger/logger"
)

// OutLogger global logger
var OutLogger *logger.Logger

func CreateLogFiles() {
	var err error
	if err = os.MkdirAll("./logs", os.ModePerm); err != nil {
		panic("Create dir './logs' error. " + err.Error())
	}
	if OutLogger, err = logger.New("logs/out.log", 1, 3, 0, logger.INFO); err != nil {
		panic("Create Outlogger file error. " + err.Error())
	}
}
