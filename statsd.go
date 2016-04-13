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

	client = statsd.NewStatsdClient(host, "goworker.")
	err := client.CreateSocket()
	if err != nil {
		log.Fatal(err)
	}
}

// Wrapper wraps goworkers and reports job duration and success/failures
func Wrapper(jobMetricName string, w func(string, ...interface{}) error) func(string, ...interface{}) error {

	// This appends a tag for sysdig
	postfix := "#job=" + jobMetricName

	return func(queue string, args ...interface{}) error {
		startTime := time.Now()
		err := w(queue, args...)
		duration := time.Since(startTime)

		client.PrecisionTiming(timeMetric+postfix, duration)
		if err == nil {
			// Increment success count
			client.Incr(successCountMetric+postfix, 1)
			client.PrecisionTiming(successTimeMetric+postfix, duration)
		} else {
			// Increment fail count
			client.Incr(failedCountMetric+postfix, 1)
			client.PrecisionTiming(failedTimeMetric+postfix, duration)
		}
		return err
	}
}
