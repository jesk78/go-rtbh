package lists

const HISTORY = "_rtbh_history"

type History struct {
}

func (hl *History) Update(address string) bool {
	key := HISTORY + ":" + address

	_, err := Redis.Incr(key).Result()
	if err != nil {
		Log.Warning("[History]: Failed to update history entry for " + address + ": " + err.Error())
		return false
	}

	return true
}

func NewHistoryList() (hl *History) {
	hl = &History{}

	return
}
