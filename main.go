package main

import (
	"fmt"
	"log"
	"net/http"
	"redirect-tool/analytics"
	"redirect-tool/healthcheck"
	"redirect-tool/http_handler"
	"redirect-tool/redis_service"
)

const apiPort = ":42069"

func main() {
	testUrls := []string{
		"https://www.youtube.com",
		"https://www.facebook.com",
		"https://www.instagram.com",
	}
	checkResults := healthcheck.CheckUrls(testUrls)
	for _, result := range checkResults {
		fmt.Println(result)
	}

	redisService, err := redis_service.NewRedisService(
		"localhost:6379",
		"",
		0,
	)

	if err != nil {
		panic(err)
	}

	manager, err := analytics.NewAnalyticsDbManager("./analytics.db")
	if err != nil {
		panic(err)
	}
	defer func(manager *analytics.DbManager) {
		err := manager.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(manager)

	err = manager.RunMigrations()
	if err != nil {
		panic(err)
	}

	httpHandler := http_handler.NewHandler(redisService, manager)
	http.HandleFunc("/shorten", httpHandler.ShortenUrlHandler)
	http.HandleFunc("/", httpHandler.RedirectHandler)

	fmt.Println("Listening on port " + apiPort)
	log.Fatal(http.ListenAndServe(apiPort, nil))
}
