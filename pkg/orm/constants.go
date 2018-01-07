package orm

import (
	"pg"

	"github.com/r3boot/go-rtbh/pkg/config"
	"github.com/r3boot/go-rtbh/pkg/logger"
)

const (
	T_ADDRESS   = "address"
	T_REASON    = "reason"
	T_DURATION  = "duration"
	T_WHITELIST = "whitelist"
	T_BLACKLIST = "blacklist"
	T_HISTORY   = "history"

	F_ID       = "id"
	F_ADDR     = "addr"
	F_FQDN     = "fqdn"
	F_REASON   = "reason"
	F_DURATION = "duration"
)

type ORM struct {
	cfg *config.Config
	log *logger.Logger
	db  *pg.DB
}

var (
	localORM *ORM
)
