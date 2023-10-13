package main

import (
    "fmt"
	"context"
    "time"

    pb  "djhuang.top/cacheserver/cache"
)


func CacheGet(client pb.CacheClient, req *pb.GetRequest) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := client.GetCache(ctx, req)
	if err != nil {
		fmt.Println("client.GetCache failed.")
	}
}

func CacheSet(client pb.CacheClient, req *pb.SetRequest) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := client.SetCache(ctx, req)
	if err != nil {
		fmt.Println("client.SetCache failed.")
	}
}

func CacheDelete(client pb.CacheClient, req *pb.DeleteRequest) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := client.DeleteCache(ctx, req)
	if err != nil {
		fmt.Println("client.DeleteCache failed.")
	}
}
