package orm

import (
	"fmt"

	"pg"
	"pg/orm"
)

func (o *ORM) Connect() error {
	localORM.db = pg.Connect(&pg.Options{
		Addr:     o.cfg.Database.Address,
		User:     o.cfg.Database.Username,
		Password: o.cfg.Database.Password,
		Database: o.cfg.Database.Name,
	})

	if o.db == nil {
		return fmt.Errorf("ORM.Connect: Failed to connect to database")
	}
	o.log.Debugf("ORM.Connect: Connected to pg://%s:***@%s/%s", o.cfg.Database.Username, o.cfg.Database.Address, o.cfg.Database.Name)

	for _, model := range []interface{}{&Address{}, &Reason{}, &Blacklist{}, &Whitelist{}, &History{}} {
		err := o.db.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			return fmt.Errorf("ORM.Connect db.CreateTable: %v", err)
		}

	}

	return nil
}
