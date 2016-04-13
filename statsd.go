package gwstatsd

import (
	"github.com/cactus/go-statsd-client/statsd"
	"log"
	"os"
	"time"
)

var client statsd.Statter

const timeMetric string = "time"
const successTimeMetric string = "success.time"
const successCountMetric string = "success.count"
const failedTimeMetric string = "fail.time"
const failedCountMetric string = "fail.count"

func init() {
	var err error

	host := os.Getenv("STATSD_HOST")
	if host == "" {
		host = "127.0.0.1:8125"
	}

	client, err = statsd.NewClient(host, "goworker")
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

		client.TimingDuration(timeMetric, duration, 1.0)
		if err == nil {
			// Increment success count
			client.Inc(successCountMetric, 1, 1.0)
			client.TimingDuration(successTimeMetric, duration, 1.0)
		} else {
			// Increment fail count
			client.Inc(failedCountMetric, 1, 1.0)
			client.TimingDuration(failedTimeMetric, duration, 1.0)
		}
		return err
	}
}
