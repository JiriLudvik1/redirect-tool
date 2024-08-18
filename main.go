package main

import (
	"fmt"
	"log"
	"net/http"
	"redirect-tool/http_handler"
	"redirect-tool/redis_service"
)

const apiPort = ":42069"

func main() {
	redisService, err := redis_service.NewRedisService(
		"localhost:6379",
		"",
		0,
	)

	if err != nil {
		panic(err)
	}

	//redirectHash, err := redisService.CreateRedirectEntry("www.facebook.com")
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("Redirect hash is: ", redirectHash)
	//
	//originalUrl, err := redisService.GetOriginalUrl(redirectHash)
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Println("OriginalUrl is: ", originalUrl)

	httpHandler := http_handler.NewHandler(redisService)
	http.HandleFunc("/shorten", httpHandler.ShortenUrlHandler)
	http.HandleFunc("/", httpHandler.RedirectHandler)

	fmt.Println("Listening on port " + apiPort)
	log.Fatal(http.ListenAndServe(apiPort, nil))
}
