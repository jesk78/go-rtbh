package proto

import (
	"errors"
	"gopkg.in/pg.v4"
)

var db_schema []string = []string{
	"CREATE TABLE IF NOT EXISTS addresses (id serial, addr text, fqdn text)",
	"CREATE TABLE IF NOT EXISTS reasons (id serial, reason text)",
	"CREATE TABLE IF NOT EXISTS durations (id serial, duration text)",
	"CREATE TABLE IF NOT EXISTS blacklists (id serial, address_id bigint, reason_id bigint, added_at timestamp, duration_id bigint)",
	"CREATE TABLE IF NOT EXISTS whitelists (id serial, addr_id bigint, description_id text)",
	"CREATE TABLE IF NOT EXISTS histories (id serial, addr_id bigint, reason_id bigint, added_at timestamp)",
}

func ConnectToPostgresql() (db *pg.DB, err error) {
	var schema_query string

	db = pg.Connect(&pg.Options{
		Addr:     Config.Database.Address,
		User:     Config.Database.Username,
		Password: Config.Database.Password,
		Database: Config.Database.Name,
	})

	if db == nil {
		err = errors.New("[ConnectToPostgresql]: Failed to connect to database")
		return
	}

	for _, schema_query = range db_schema {
		_, err = db.Exec(schema_query)
		if err != nil {
			err = errors.New("[ConnectToPostgresql]: Failed to update DB schema: " + err.Error())
			return
		}
	}

	return
}
