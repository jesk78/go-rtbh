package orm

import (
	"github.com/r3boot/go-rtbh/lib/config"
	"github.com/r3boot/rlib/logger"
	"gopkg.in/pg.v4"
)

const MYNAME string = "ORM"

var Config *config.Config
var Log logger.Log

var db *pg.DB

type ORM struct {
}

func Setup(l logger.Log, c *config.Config) (err error) {
	Log = l
	Config = c

	Log.Debug(MYNAME + ": Module initialized")
	return
}

func New() *ORM {
	var orm *ORM

	orm = &ORM{}
	orm.Connect()

	return orm
}
