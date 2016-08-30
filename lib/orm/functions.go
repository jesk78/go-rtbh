package orm

import (
	"gopkg.in/pg.v4"
)

func (orm *ORM) Connect() (err error) {
	var schema_query string

	orm.Db = pg.Connect(&pg.Options{
		Addr:     Config.Database.Address,
		User:     Config.Database.Username,
		Password: Config.Database.Password,
		Database: Config.Database.Name,
	})

	if db == nil {
		err = errors.New(MYNAME + ": Failed to connect to database")
		return
	}

	for _, schema_query = range database_schema {
		if _, err = orm.Db.Exec(schema_query); err != nil {
			err = errors.New(MYNAME + ": Failed to update DB schema: " + err.Error())
			return
		}
	}

	return
}
