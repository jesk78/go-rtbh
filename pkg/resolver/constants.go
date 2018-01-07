package resolver

import (
	"sync"
	"time"

	"github.com/r3boot/go-rtbh/pkg/config"
	"github.com/r3boot/go-rtbh/pkg/logger"
	"github.com/r3boot/go-rtbh/pkg/orm"
)

const (
	FQDN_TO_LOOKUP = "not-yet-resolved"
	FQDN_NXHOST    = "nxhost"

	MAX_SAMPLES        = 100
	MAX_SLEEP_INTERVAL = 5000
)

type Resolver struct {
	cfg      *config.Config
	log      *logger.Logger
	orm      *orm.ORM
	Interval time.Duration
	Control  chan int
	Done     chan bool
	cache    map[string]string
	mutex    *sync.Mutex
}
