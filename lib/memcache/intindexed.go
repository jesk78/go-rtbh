package memcache

func (cache *IntCache) Has(key int64) bool {
	var cached_key int64

	for cached_key, _ = range cache.cache {
		if cached_key == key {
			return true
		}
	}

	return false
}

func (cache *IntCache) Add(key int64, entry interface{}) {
	cache.mutex.Lock()
	cache.cache[key] = entry
	cache.mutex.Unlock()
}

func (cache *IntCache) AddUnsafe(key int64, entry interface{}) {
	cache.cache[key] = entry
}

func (cache *IntCache) Remove(key int64) {
	if ok := cache.Has(key); !ok {
		return
	}

	cache.mutex.Lock()
	delete(cache.cache, key)
	cache.mutex.Unlock()
}

func (cache *IntCache) GetAll() map[int64]interface{} {
	return cache.cache
}

func (cache *IntCache) Get(key int64) interface{} {
	return cache.cache[key]
}
