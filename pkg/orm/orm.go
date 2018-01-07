package orm

import (
	"fmt"

	"github.com/r3boot/go-rtbh/pkg/config"
	"github.com/r3boot/go-rtbh/pkg/logger"
)

func New(l *logger.Logger, c *config.Config) (*ORM, error) {

	localORM = &ORM{
		log: l,
		cfg: c,
	}
	err := localORM.Connect()
	if err != nil {
		return nil, fmt.Errorf("localORM.New: %v", err)
	}

	return localORM, nil
}
