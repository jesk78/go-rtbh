package config

// Default configuration options
const D_APICFGFILE string = "go-rtbhapi.yml"
const D_CFGFILE string = "go-rtbh.yml"
const D_DEBUG bool = false
const D_TIMESTAMP bool = false

// Application defaults
const REDIS_D_ADDR string = "localhost:6379"
const AMQP_D_ADDR string = "localhost:5672"
const AMQP_D_USER string = "go-rtbh"
const AMQP_D_PASS string = "go-rtbh"
const AMQP_D_EXCHANGE string = "amqp-input"

// Channel buffer sizes
const D_SIGNAL_BUFSIZE int = 16
const D_CONTROL_BUFSIZE int = 1
const D_DONE_BUFSIZE int = 1
const D_REDIS_BUFSIZE int = 32
const D_AMQP_BUFSIZE int = 32
const D_INPUT_BUFSIZE int = 64

// Goroutine control signals
const CTL_SHUTDOWN int = 0

// Various files used
const TMPL_BIRD string = "/usr/share/go-rtbh/bird.conf.template"
