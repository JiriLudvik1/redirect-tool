package healthcheck

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type CheckResult struct {
	URL         string
	statusCode  int
	error       error
	elapsedTime time.Duration
}

func CheckUrls(urls []string) []CheckResult {
	resultsChannel := make(chan CheckResult)

	// Start a goroutine to close the channel after all checks are done
	go func() {
		for _, url := range urls {
			go checkUrl(url, resultsChannel)
		}
	}()

	// The number of URLs to check
	expectedResultsCount := len(urls)

	// Collect all results
	results := make([]CheckResult, len(urls))
	for i := 0; i < expectedResultsCount; i++ {
		results[i] = <-resultsChannel
	}

	// Close the channel after receiving all results
	close(resultsChannel)
	return results
}

func checkUrl(url string, resultsChannel chan<- CheckResult) {
	startTime := time.Now()

	resp, err := http.Get(url)
	elapsedTime := time.Since(startTime)

	if err != nil {
		resultsChannel <- CheckResult{url, 0, err, elapsedTime}
		return
	}
	if resp == nil {
		resultsChannel <- CheckResult{url, 500, nil, elapsedTime}
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	result := CheckResult{URL: url, statusCode: resp.StatusCode, error: nil, elapsedTime: elapsedTime}
	resultsChannel <- result
}
