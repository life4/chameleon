package chameleon

import (
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	cache "github.com/victorspringer/http-cache"
	"github.com/victorspringer/http-cache/adapter/memory"
)

type Cache struct {
	client *cache.Client
}

func NewCache(capacity int) (*Cache, error) {
	if capacity == 0 {
		return nil, nil
	}
	memcached, err := memory.NewAdapter(
		memory.AdapterWithAlgorithm(memory.LRU),
		memory.AdapterWithCapacity(capacity),
	)
	if err != nil {
		return nil, err
	}

	client, err := cache.NewClient(
		cache.ClientWithAdapter(memcached),
		cache.ClientWithTTL(10*time.Minute),
	)
	if err != nil {
		return nil, err
	}
	return &Cache{client: client}, nil
}

func (c *Cache) Wrap(f httprouter.Handle) httprouter.Handle {
	if c == nil {
		return f
	}
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		h1 := func(w http.ResponseWriter, r *http.Request) {
			f(w, r, ps)
		}
		h2 := c.client.Middleware(http.HandlerFunc(h1))
		h2.ServeHTTP(w, r)
	}
}
