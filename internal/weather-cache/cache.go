package weather_cache

import (
	"github.com/Back1ng/openmeteo/internal/entity"
	"sync"
	"time"
)

type Cache struct {
	Weather entity.Weather

	ttl time.Time
	mu  *sync.RWMutex
}

func New() *Cache {
	return &Cache{
		ttl: time.Now(),
		mu:  &sync.RWMutex{},
	}
}

func (c *Cache) Get() (*entity.Weather, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.IsValid() {
		return nil, false
	}

	return &c.Weather, true
}

func (c *Cache) Store(weather entity.Weather) {
	c.mu.Lock()
	c.Weather = weather
	c.ttl = time.Now().Add(time.Minute)
	c.mu.Unlock()
}

func (c *Cache) IsValid() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return time.Now().Before(c.ttl) && c.Weather.Temp != float32(0)
}
