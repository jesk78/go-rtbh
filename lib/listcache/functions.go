package listcache

func (cache *Cache) Has(addr string) bool {
	var cached_addr string

	for cached_addr, _ = range cache.cache {
		if cached_addr == addr {
			return true
		}
	}

	return false
}

func (cache *Cache) Add(addr string, entry interface{}) {
	cache.mutex.Lock()
	cache.cache[addr] = entry
	cache.mutex.Unlock()
}

func (cache *Cache) Remove(addr string) {
	if ok := cache.Has(addr); !ok {
		return
	}

	cache.mutex.Lock()
	delete(cache.cache, addr)
	cache.mutex.Unlock()
}

func (cache *Cache) GetAll() map[string]interface{} {
	return cache.cache
}
