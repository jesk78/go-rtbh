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
	"CREATE TABLE IF NOT EXISTS addresses (id serial, addr text, fqdn text)",
	"CREATE TABLE IF NOT EXISTS reasons (id serial, reason text)",
	"CREATE TABLE IF NOT EXISTS durations (id serial, duration text)",
	"CREATE TABLE IF NOT EXISTS blacklists (id serial, address_id bigint, reason_id bigint, added_at timestamp, duration_id bigint)",
	"CREATE TABLE IF NOT EXISTS whitelists (id serial, addr_id bigint, description_id text)",
	"CREATE TABLE IF NOT EXISTS histories (id serial, addr_id bigint, reason_id bigint, added_at timestamp)",
}
