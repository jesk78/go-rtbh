package memcache

func (cache *StringCache) Has(key string) bool {
	var cached_key string

	for cached_key, _ = range cache.cache {
		if cached_key == key {
			return true
		}
	}

	return false
}

func (cache *StringCache) Add(key string, entry interface{}) {
	cache.mutex.Lock()
	cache.cache[key] = entry
	cache.mutex.Unlock()
}

func (cache *StringCache) AddUnsafe(key string, entry interface{}) {
	cache.cache[key] = entry
}

func (cache *StringCache) Remove(key string) {
	if ok := cache.Has(key); !ok {
		return
	}

	cache.mutex.Lock()
	delete(cache.cache, key)
	cache.mutex.Unlock()
}

func (cache *StringCache) GetAll() map[string]interface{} {
	return cache.cache
}

func (cache *StringCache) Get(key string) interface{} {
	return cache.cache[key]
}
