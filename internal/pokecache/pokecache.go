package pokecache

import (
  "time" 
  "sync"
)

type cacheEntry struct{
  createdAt time.Time
  val []byte
}

type Cache struct{
  mu sync.Mutex
  cache map[string]cacheEntry
}

func (c Cache) Add(key string, val []byte){
  c.mu.Lock()
  _, ok := c.cache[key]
  if !ok{
    var currentCacheEntry cacheEntry
    currentCacheEntry.createdAt = time.Now()
    currentCacheEntry.val = val
    c.cache[key] = currentCacheEntry
  }
  c.mu.Unlock()
}

func (c Cache) Get(key string) ([]byte, bool){
  c.mu.Lock()
  entry, ok := c.cache[key]
  if !ok{
    return []byte{}, false
  }
  return entry.val, true
}

func reapLoop(c Cache, interval time.Duration){
  ticker := time.NewTicker(1 * time.Millisecond)
  go func()  {
    for{
      select{
      case <-ticker.C:
        for m, cache := range c.cache{
          createdTime := cache.createdAt
          compareTime := time.Now()
          //createdTime.Add(interval * time.Second)
          if compareTime.Sub(createdTime) >= interval{
            delete(c.cache, m)
          }
        }

    }
    }
  }()

}

func NewCache(interval time.Duration) Cache{
  cache := make(map[string]cacheEntry)
  mu := sync.Mutex{}
  currentCache := Cache{mu, cache}
  go reapLoop(currentCache, interval)
  return currentCache
} 
