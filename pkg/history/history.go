package history

import (
	"github.com/r3boot/go-rtbh/pkg/config"
	"github.com/r3boot/go-rtbh/pkg/logger"
)

func NewHistory(l *logger.Logger, c *config.Config) *History {
	log = l
	cfg = c

	history := &History{}

	log.Debugf("History: Module initialized")

	return history
}
