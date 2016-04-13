package gwstatsd

import (
	"github.com/quipo/statsd"
	"log"
	"os"
	"time"
)

var client *statsd.StatsdClient

const timeMetric string = "time"
const successTimeMetric string = "success.time"
const successCountMetric string = "success.count"
const failedTimeMetric string = "fail.time"
const failedCountMetric string = "fail.count"

func init() {
	host := os.Getenv("STATSD_HOST")
	if host == "" {
		host = "127.0.0.1:8125"
	}

	client = statsd.NewStatsdClient(host, "goworker")
	err := client.CreateSocket()
	if err != nil {
		log.Fatal(err)
	}
}

// Wrapper wraps goworkers and reports job duration and success/failures
func Wrapper(w func(string, ...interface{}) error) func(string, ...interface{}) error {

	return func(queue string, args ...interface{}) error {
		startTime := time.Now()
		err := w(queue, args...)
		duration := time.Since(startTime)

		client.PrecisionTiming(timeMetric, duration)
		if err == nil {
			// Increment success count
			client.Incr(successCountMetric, 1)
			client.PrecisionTiming(successTimeMetric, duration)
		} else {
			// Increment fail count
			client.Incr(failedCountMetric, 1)
			client.PrecisionTiming(failedTimeMetric, duration)
		}
		return err
	}
}
