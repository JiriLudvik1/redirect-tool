package http_handler

import "time"

func RateLimiter(limiterChan chan struct{}) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		for i := 0; i < 10; i++ {
			select {
			case limiterChan <- struct{}{}:
			default:
			}
		}
	}
}
