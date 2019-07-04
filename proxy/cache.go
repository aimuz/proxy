package proxy

import (
	"sync"
	"time"
)

type Cache struct {
	ExecAt time.Time
	Info   *Info
	List   *List
}

var caches = make(map[string]*Cache)
var clock = sync.Mutex{}

func getCache(_path string, version string) (*Cache, bool) {
	clock.Lock()
	defer clock.Unlock()

	v, ok := caches[_path+"@"+version]
	if ok {
		if (version != "latest") ||
			(version == "latest" && time.Now().Sub(v.ExecAt).Seconds() < 30) {
			return v, ok
		}
	}
	return nil, false
}

func setCache(key string, cache *Cache) {
	clock.Lock()
	defer clock.Unlock()
	caches[key] = cache
}
