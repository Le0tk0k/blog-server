package log

import (
	"github.com/Le0tk0k/blog-server/config"

	"github.com/labstack/gommon/log"
)

// New はロガーを生成する
func New() *log.Logger {
	logger := log.New("application")
	logger.SetLevel(log.INFO)
	if config.IsLocal() {
		logger.SetLevel(log.DEBUG)
	}
	return logger
}
