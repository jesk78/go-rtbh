package orm

const (
	T_ADDRESS   string = "address"
	T_REASON    string = "reason"
	T_DURATION  string = "duration"
	T_WHITELIST string = "whitelist"
	T_BLACKLIST string = "blacklist"
	T_HISTORY   string = "history"

	F_ID       string = "id"
	F_ADDR     string = "addr"
	F_FQDN     string = "fqdn"
	F_REASON   string = "reason"
	F_DURATION string = "duration"
)

var database_schema []string = []string{
	"CREATE TABLE IF NOT EXISTS addresses (id SERIAL PRIMARY KEY, addr TEXT UNIQUE NOT NULL, fqdn TEXT NOT NULL)",
	"CREATE TABLE IF NOT EXISTS reasons (id SERIAL PRIMARY KEY, reason TEXT UNIQUE NOT NULL)",
	"CREATE TABLE IF NOT EXISTS durations (id SERIAL PRIMARY KEY, duration TEXT UNIQUE NOT NULL)",
	"CREATE TABLE IF NOT EXISTS blacklists (id SERIAL PRIMARY KEY, addr_id BIGINT NOT NULL REFERENCES addresses(id), reason_id BIGINT NOT NULL REFERENCES reasons(id), added_at TIMESTAMP NOT NULL, duration_id BIGINT NOT NULL REFERENCES durations(id))",
	"CREATE TABLE IF NOT EXISTS whitelists (id SERIAL PRIMARY KEY, addr_id BIGINT NOT NULL REFERENCES addresses(id), description TEXT NOT NULL)",
	"CREATE TABLE IF NOT EXISTS histories (id SERIAL PRIMARY KEY, addr_id BIGINT NOT NULL REFERENCES addresses(id), reason_id BIGINT NOT NULL REFERENCES reasons(id), added_at timestamp)",
}
