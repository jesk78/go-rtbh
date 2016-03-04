package lists

import (
	"strings"
)

const WHITELIST = "_rtbh_whitelist"

type Whitelist struct {
}

func (wl *Whitelist) Add(addr string, description string) bool {
	_, err := Redis.HMSet(WHITELIST, addr, description).Result()
	if err != nil {
		Log.Warning("[Whitelist]: Failed to add " + addr + ": " + err.Error())
		return false
	}

	return true
}

func (wl *Whitelist) Count() int64 {
	var count int64
	var err error

	if count, err = Redis.HLen(WHITELIST).Result(); err != nil {
		Log.Warning("[Whitelist.Count()]: Redis.HLen(): " + err.Error())
		return -1
	}

	return count
}

func (wl *Whitelist) Listed(address string) bool {
	result := Redis.HMGet(WHITELIST, address).Val()[0]
	return result != nil
}

func (wl *Whitelist) GetAll() []string {
	var result []string

	for _, addr := range Redis.HKeys(WHITELIST).Val() {
		// Skip IPv6 addresses for now
		if strings.Contains(addr, ":") {
			continue
		}

		result = append(result, addr)
	}

	return result
}

func NewWhitelist() (wl *Whitelist) {
	wl = &Whitelist{}

	for _, entry := range Config.Whitelist {
		Log.Debug("[Whitelist]: Adding " + entry.Address + " to the whitelist")
		wl.Add(entry.Address, entry.Description)
	}

	return
}
