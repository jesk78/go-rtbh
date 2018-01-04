package orm

import (
	"fmt"

	"gopkg.in/pg.v4"
)

func (orm *ORM) Connect() error {
	db = pg.Connect(&pg.Options{
		Addr:     cfg.Database.Address,
		User:     cfg.Database.Username,
		Password: cfg.Database.Password,
		Database: cfg.Database.Name,
	})

	if db == nil {
		return fmt.Errorf("ORM.Connect: Failed to connect to database")
	}
	log.Debugf("ORM.Connect: Connected to pg://%s:***@%s/%s", cfg.Database.Username, cfg.Database.Address, cfg.Database.Name)

	for _, schema_query := range databaseSchema {
		_, err := db.Exec(schema_query)
		if err != nil {
			log.Debugf(schema_query)
			return fmt.Errorf("ORM.Connect db.Exec: %v", err)
		}
	}

	return nil
}
