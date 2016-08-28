package orm

import (
	"github.com/r3boot/rlib/logger"
	"gopkg.in/pg.v4"
)

var Log logger.Log
var Db *pg.DB

func Setup(l logger.Log, db *pg.DB) (err error) {
	Log = l
	Db = db

	return
}
