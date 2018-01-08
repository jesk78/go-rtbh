package config

const (
	// Default configuration options
	D_APICFGFILE = "go-rtbhapi.yml"
	D_CFGFILE    = "go-rtbh.yml"
	D_DEBUG      = false
	D_TIMESTAMP  = false

	// Application defaults
	REDIS_D_ADDR    = "localhost:6379"
	AMQP_D_ADDR     = "localhost:5672"
	AMQP_D_USER     = "go-rtbh"
	AMQP_D_PASS     = "go-rtbh"
	AMQP_D_EXCHANGE = "amqp-input"

	// Channel buffer sizes
	D_SIGNAL_BUFSIZE  = 16
	D_CONTROL_BUFSIZE = 1
	D_DONE_BUFSIZE    = 1
	D_REDIS_BUFSIZE   = 32

	D_AMQP_BUFSIZE  = 32
	D_INPUT_BUFSIZE = 64

	// Goroutine control signals
	CTL_SHUTDOWN = 0

	// Various files used
	TMPL_BIRD = "/usr/share/go-rtbh/bird.conf.template"
)

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

type ESConfig struct {
	Address string `yaml:"address"`
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
	ES        ESConfig          `yaml:"elasticsearch"`
	Redis     RedisConfig       `yaml:"redis"`
	Database  PostgresConfig    `yaml:"postgresql"`
	BGP       BGPConfig         `yaml:"bgp"`
	Whitelist []WhitelistConfig `yaml:"whitelist"`
	Blacklist []BlacklistConfig `yaml:"blacklist"`
	Ruleset   []string          `yaml:"ruleset"`
}
