package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/charmbracelet/log"
	"github.com/redis/go-redis/v9"
)

type Cache struct {
	Db     *redis.Client
	Ctx    context.Context
	Prefix string
}

// NewCache creates a new database logger
func NewCache(inDb *redis.Client, inPrefix string) *Cache {
	return &Cache{
		Db:     inDb,
		Ctx:    context.Background(),
		Prefix: inPrefix,
	}
}

func (cache *Cache) Set(inKey string, inValue any, inTime time.Duration) {
	jsonString, _ := json.Marshal(inValue)
	err := cache.Db.Set(cache.Ctx, cache.Prefix+inKey, string(jsonString), inTime).Err()
	if err != nil {
		log.Error(err)
	}
}

func formatData(inValue string) any {
	var opValue any
	err := json.Unmarshal([]byte(inValue), opValue)
	if err != nil {
		log.Error(err)
		return nil
	}
	return opValue
}
func (cache *Cache) Get(inKey string) any {
	cacheData, err := cache.Db.Get(cache.Ctx, cache.Prefix+inKey).Result()
	if err == nil {
		return formatData(cacheData)
	}
	return nil
}

func (cache *Cache) Del(inKey string) {
	cache.Db.Del(cache.Ctx, cache.Prefix+inKey).Result()
}

func (cache *Cache) Dels(inKeys ...string) {
	for i := 0; i < len(inKeys); i += 1 {
		inKeys[i] = cache.Prefix + inKeys[i]
	}
	cache.Db.Del(cache.Ctx, inKeys...).Result()
}

func (cache *Cache) HSet(inKey string, inProperty, inValue any) {
	jsonString, _ := json.Marshal(inValue)
	err := cache.Db.HSet(cache.Ctx, cache.Prefix+inKey, inProperty, jsonString).Err()
	if err != nil {
		log.Error(err)
	}
}

func (cache *Cache) HGet(inKey string, inProperty string) any {
	cacheData, err := cache.Db.HGet(cache.Ctx, cache.Prefix+inKey, inProperty).Result()
	if err == nil {
		return formatData(cacheData)
	}
	return nil
}

func (cache *Cache) HDel(inKey string, inProperty string) {
	cache.Db.HDel(cache.Ctx, cache.Prefix+inKey, inProperty).Result()
}
