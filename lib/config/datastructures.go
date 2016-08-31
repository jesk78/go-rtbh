package config

// Configuration structs
type ResolverConfig struct {
	Enabled           bool   `yaml:"enable"`
	LookupMaxInterval string `yaml:"max_interval"`
}

type GeneralConfig struct {
	NumWorkers     int            `yaml:"workers"`
	ReaperInterval string         `yaml:"reaper_interval"`
	Resolver       ResolverConfig `yaml:"resolver"`
}

type ApiConfig struct {
	BindIp    string `yaml:"bindip"`
	BindPort  string `yaml:"bindport"`
	Resources string `yaml:"resources"`
}

type AmqpConfig struct {
	Address  string `yaml:"address"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Exchange string `yaml:"exchange"`
}

type RedisConfig struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
	Database int64  `yaml:"database"`
}

type PostgresConfig struct {
	Address  string `yaml:"address"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type BGPPeer struct {
	Name    string `yaml:"name"`
	Address string `yaml:"address"`
	Asnum   string `yaml:"asnum"`
}

type BGPConfig struct {
	Asnum     string    `yaml:"asnum"`
	RouterId  string    `yaml:"routerid"`
	NextHop   string    `yaml:"nexthop"`
	NextHopV6 string    `yaml:"nexthopv6"`
	Community string    `yaml:"community"`
	LocalPref int       `yaml:localpref"`
	Peers     []BGPPeer `yaml:"peers"`
}

type WhitelistConfig struct {
	Address     string `yaml:"address"`
	Description string `yaml:"description"`
}

type BlacklistConfig struct {
	Address string `yaml:"address"`
	Reason  string `yaml:"reason"`
}

type Config struct {
	General   GeneralConfig     `yaml:"general"`
	Api       ApiConfig         `yaml:"api"`
	Amqp      AmqpConfig        `yaml:"amqp"`
	Redis     RedisConfig       `yaml:"redis"`
	Database  PostgresConfig    `yaml:"postgresql"`
	BGP       BGPConfig         `yaml:"bgp"`
	Whitelist []WhitelistConfig `yaml:"whitelist"`
	Blacklist []BlacklistConfig `yaml:"blacklist"`
	Ruleset   []string          `yaml:"ruleset"`
}
