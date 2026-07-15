package main

import (
	"context"
	"fmt"
	"log"

	"github.com/you/sharing-vision-backend-v2/internal/cache"
	"github.com/you/sharing-vision-backend-v2/internal/config"
)

func main() {
	config.Load()
	client, err := cache.ConnectRedis()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	ctx := context.Background()
	prefix := "sv:articles:*"
	keys, err := client.Keys(ctx, prefix).Result()
	if err != nil {
		log.Fatal(err)
	}
	if len(keys) == 0 {
		fmt.Println("no cache keys")
		return
	}
	if err := client.Del(ctx, keys...).Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("flushed %d cache keys\n", len(keys))
}
