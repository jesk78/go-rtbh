package history

import (
	"github.com/r3boot/go-rtbh/pkg/config"
	"github.com/r3boot/go-rtbh/pkg/logger"
	"github.com/r3boot/go-rtbh/pkg/orm"
)

type History struct {
	cfg *config.Config
	log *logger.Logger
	orm *orm.ORM
}
