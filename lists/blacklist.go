package lists

import (
	"strings"
)

const BLACKLIST = "_rtbh_blacklist"

type Blacklist struct {
}

func (bl *Blacklist) Add(address string, reason string) bool {
	_, err := Redis.HMSet(BLACKLIST, address, reason).Result()
	if err != nil {
		Log.Warning("[Blacklist]: Failed to add " + address + ": " + err.Error())
		return false
	}
	return true
}

func (bl *Blacklist) Listed(address string) bool {
	result := Redis.HMGet(BLACKLIST, address).Val()[0]
	return result != nil
}

func (bl *Blacklist) GetAll() []string {
	var result []string

	for _, addr := range Redis.HKeys(BLACKLIST).Val() {
		// Skip IPv6 addresses for now
		if strings.Contains(addr, ":") {
			continue
		}
		result = append(result, addr)
	}

	return result
}

func NewBlacklist() (bl *Blacklist) {
	bl = &Blacklist{}

	for _, entry := range Config.Blacklist {
		Log.Debug("[Blacklist]: Adding " + entry.Address + " to the blacklist: " + entry.Reason)
		bl.Add(entry.Address, entry.Reason)
	}

	return
}
