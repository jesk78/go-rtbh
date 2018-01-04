package orm

import (
	"fmt"

	"github.com/r3boot/go-rtbh/pkg/config"
	"github.com/r3boot/go-rtbh/pkg/logger"
	"github.com/r3boot/go-rtbh/pkg/memcache"
	"gopkg.in/pg.v4"
)

var db *pg.DB

func New(l *logger.Logger, c *config.Config) (*ORM, error) {
	log = l
	cfg = c

	orm := &ORM{}
	err := orm.Connect()
	if err != nil {
		return nil, fmt.Errorf("NewORM: %v", err)
	}

	// Initialize the various caches
	cacheAddrOnAddress = memcache.NewStringIndexed()
	cacheAddrOnId = memcache.NewIntIndexed()
	err = WarmupAddressCaches()
	if err != nil {
		return nil, fmt.Errorf("NewORM: %v", err)
	}

	cacheReasonOnAddress = memcache.NewStringIndexed()
	cacheReasonOnId = memcache.NewIntIndexed()
	err = WarmupReasonCaches()
	if err != nil {
		return nil, fmt.Errorf("NewORM: %v", err)
	}

	cacheBlacklistOnAddress = memcache.NewStringIndexed()
	cacheBlacklistOnId = memcache.NewIntIndexed()
	err = WarmupBlacklistCaches()
	if err != nil {
		return nil, fmt.Errorf("NewORM: %v", err)
	}

	cacheHistory = memcache.NewStringIndexed()
	// TODO: Finish history cache

	return orm, nil
}
