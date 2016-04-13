# goworker-statsd
goworker wrapper that reports statsd success, failure, and job duration
times.

## Usage
```go
goworker.Register("MyClass", gwstatsd.Wrapper("myclass", myWorker))
```
