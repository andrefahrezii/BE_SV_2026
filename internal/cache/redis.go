package cache

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/you/sharing-vision-backend-v2/internal/config"
	"github.com/you/sharing-vision-backend-v2/internal/model"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	Client *redis.Client
}

func NewCache(client *redis.Client) *Cache {
	return &Cache{Client: client}
}

func ConnectRedis() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Conf.Redis.Addr,
		Password: config.Conf.Redis.Password,
		DB:       config.Conf.Redis.DB,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	config.Log.Info("redis connected")
	return rdb, nil
}

func (c *Cache) buildKeyPrefix() string {
	return "sv:articles"
}

func (c *Cache) BuildListKey(q model.PostListQuery) string {
	return c.buildKeyPrefix() + ":page:" + strconv.Itoa(q.Offset) + ":limit:" + strconv.Itoa(q.Limit) + ":category:" + q.Category + ":q:" + q.Q + ":status:" + q.Status
}

func (c *Cache) GetList(q model.PostListQuery) ([]model.PostListItem, int, bool, error) {
	if q.Limit <= 0 {
		q.Limit = 10
	}
	key := c.BuildListKey(q)
	val, err := c.Client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return nil, 0, false, nil
	}
	if err != nil {
		return nil, 0, false, err
	}
	var result struct {
		Items []model.PostListItem `json:"items"`
		Total int                  `json:"total"`
	}
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		return nil, 0, false, err
	}
	return result.Items, result.Total, true, nil
}

func (c *Cache) SetList(q model.PostListQuery, items []model.PostListItem, total int) {
	if q.Limit <= 0 {
		q.Limit = 10
	}
	key := c.BuildListKey(q)
	result := struct {
		Items []model.PostListItem `json:"items"`
		Total int                  `json:"total"`
	}{Items: items, Total: total}
	b, _ := json.Marshal(result)
	c.Client.Set(context.Background(), key, b, 5*time.Minute)
}

func (c *Cache) InvalidateList() error {
	ctx := context.Background()
	keys, err := c.Client.Keys(ctx, c.buildKeyPrefix()+":*").Result()
	if err != nil {
		return err
	}
	for _, k := range keys {
		c.Client.Del(ctx, k)
	}
	return nil
}
