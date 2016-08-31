package orm

import (
	"errors"
	"gopkg.in/pg.v4"
)

func (orm *ORM) Connect() (err error) {
	var schema_query string

	db = pg.Connect(&pg.Options{
		Addr:     Config.Database.Address,
		User:     Config.Database.Username,
		Password: Config.Database.Password,
		Database: Config.Database.Name,
	})

	if db == nil {
		err = errors.New(MYNAME + ": Failed to connect to database")
		return
	}
	Log.Debug(MYNAME + ": Connected to pg://" + Config.Database.Username + ":***@" + Config.Database.Address + "/" + Config.Database.Name)

	for _, schema_query = range database_schema {
		if _, err = db.Exec(schema_query); err != nil {
			err = errors.New(MYNAME + ": Failed to update DB schema: " + err.Error())
			return
		}
	}

	return
}
