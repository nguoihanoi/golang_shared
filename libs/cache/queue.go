package cache

import "encoding/json"

type Queue struct {
	cache *Cache
	name  string
}

func (cache *Cache) NewQueue(inName string) *Queue {
	return &Queue{
		cache: cache,
		name:  inName,
	}
}
func (q *Queue) Push(inValue any) {
	jsonString, _ := json.Marshal(inValue)
	q.cache.Db.RPush(q.cache.Ctx, q.name, string(jsonString))
}

func (q *Queue) Pop() any {
	cacheData, err := q.cache.Db.LPop(q.cache.Ctx, q.name).Result()
	if err == nil {
		return formatData(cacheData)
	}
	return nil
}
func (q *Queue) Count() int64 {
	cacheData, err := q.cache.Db.LLen(q.cache.Ctx, q.name).Result()
	if err == nil {
		return cacheData
	}
	return 0
}
func (q *Queue) Gets() (output []any) {
	cacheData, err := q.cache.Db.LRange(q.cache.Ctx, q.name, 0, -1).Result()
	if err == nil {
		for i := range cacheData {
			output = append(output, formatData(cacheData[i]))
		}
	}
	return output
}
