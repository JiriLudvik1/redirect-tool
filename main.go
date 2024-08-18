package main

import (
	"context"
	"fmt"
	"redirect-tool/redis_service"
)

var ctx = context.Background()

func main() {
	redisService, err := redis_service.NewRedisService(
		"localhost:6379",
		"",
		0,
	)

	if err != nil {
		panic(err)
	}

	redirectHash, err := redisService.CreateRedirectEntry("www.facebook.com")
	if err != nil {
		panic(err)
	}
	fmt.Println("Redirect hash is: ", redirectHash)

	originalUrl, err := redisService.GetOriginalUrl(redirectHash)
	if err != nil {
		panic(err)
	}

	fmt.Println("OriginalUrl is: ", originalUrl)
}
