package orm

import (
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
}

var (
	cfg *config.Config
	log *logger.Logger
)

var databaseSchema []string = []string{
	"CREATE TABLE IF NOT EXISTS addresses (id SERIAL PRIMARY KEY, addr TEXT UNIQUE NOT NULL, fqdn TEXT NOT NULL)",
	"CREATE TABLE IF NOT EXISTS reasons (id SERIAL PRIMARY KEY, reason TEXT UNIQUE NOT NULL)",
	"CREATE TABLE IF NOT EXISTS blacklists (id SERIAL PRIMARY KEY, addr_id BIGINT NOT NULL REFERENCES addresses(id), reason_id BIGINT NOT NULL REFERENCES reasons(id), added_at TIMESTAMP NOT NULL, expire_on TIMESTAMP NOT NULL)",
	"CREATE TABLE IF NOT EXISTS whitelists (id SERIAL PRIMARY KEY, addr_id BIGINT NOT NULL REFERENCES addresses(id), description TEXT NOT NULL)",
	"CREATE TABLE IF NOT EXISTS histories (id SERIAL PRIMARY KEY, addr_id BIGINT NOT NULL REFERENCES addresses(id), reason_id BIGINT NOT NULL REFERENCES reasons(id), added_at TIMESTAMP)",
}
