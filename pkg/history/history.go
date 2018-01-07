package history

import (
	"github.com/r3boot/go-rtbh/pkg/config"
	"github.com/r3boot/go-rtbh/pkg/logger"
	"github.com/r3boot/go-rtbh/pkg/orm"
)

func New(l *logger.Logger, c *config.Config, o *orm.ORM) *History {
	history := &History{
		log: l,
		cfg: c,
		orm: o,
	}

	history.log.Debugf("History: Module initialized")

	return history
}
