# challenge_fr

# README
This is an improved and fixed script for checking endpoint's availability based on input config file.

# How to run original script
add go.mod with 
```
 go mod init github.com/vovax3m/challenge_fr
 go mod tidy
```
# How to run raw script
```
 go run main.go path/to/sample.yaml
```
# How to build binary
  ```
  cd golang
  go build .
  ```
# How to execute binary
```
  ./challenge_fr path/to/sample.yaml
```
# How to build docker image
```
docker build . -t challenge_fr:1
```
# How to run dockerized app
```
//configs could be baked into the image and app executed without addition configs, but it's a bad practice, so mounting from host
// -e LOG_LEVEL is optional to enable debug informaiton
docker run -ti -e LOG_LEVEL=DEBUG -p'9090:9090' -v '/path/to/sample.yaml:/go/sample.yaml:ro' challenge_fr:1 /go/sample.yaml
```
# Feature flags
```
export LOG_LEVEL=DEBUG // to enable detailed information
unset LOG_LEVEL // to show only report and errors

export LOG_FILE=/path/log.txt //to enable writing logs to the file
unset LOG_FILE // to send  logs only to Stdout/Stderr
```


# How to get metrics
```
curl localhost:9090/metrics
```
metrics:
```
# HELP request_latency_miliseconds_last Last request latency per domain
# TYPE request_latency_miliseconds_last gauge
request_latency_miliseconds_last{domain="dev-sre-take-home-exercise-rubric.us-east-1.recruiting-public.fetchrewards.com"} 62
request_latency_miliseconds_last{domain="example.com"} 69
request_latency_miliseconds_last{domain="httpstat.us"} 330

# HELP request_latency_miliseconds HTTP request latency in milliseconds
# TYPE request_latency_miliseconds histogram
request_latency_miliseconds_bucket{domain="dev-sre-take-home-exercise-rubric.us-east-1.recruiting-public.fetchrewards.com",le="10"} 0
request_latency_miliseconds_bucket{domain="dev-sre-take-home-exercise-rubric.us-east-1.recruiting-public.fetchrewards.com",le="50"} 0
request_latency_miliseconds_bucket{domain="dev-sre-take-home-exercise-rubric.us-east-1.recruiting-public.fetchrewards.com",le="100"} 5

# HELP domain_availability_percentage Percentage of domains availability
# TYPE domain_availability_percentage gauge
domain_availability_percentage{domain="dev-sre-take-home-exercise-rubric.us-east-1.recruiting-public.fetchrewards.com"} 15
domain_availability_percentage{domain="example.com"} 100
domain_availability_percentage{domain="httpstat.us"} 17

# HELP uptime_total The total number of service uptime
# TYPE uptime_total counter
uptime_total 31

```


# Found issues with original code:

- go.mod not exists
- main.go:87 used fmt instead log
- port ignore not implemented
- latency check not implemented
- request Body took whole endpoint not endpoint.Body
- timing not implemented
- no comments or additional messages
- no input config validation

# Troubleshooting steps:
- Make application run, add go.mod
- Add debugging option to check required features
- Added timestamp for accurate time tracking in logs
- Play with sample configuration to find gaps in implementation
- Test each requirement and fix/implement in the code
- Split logTimer and checkHealth to goroutines for independent execution 
- Wrap to docker
- Update readme file
- Add config entry validation
- Moved hardcoded timers and thresholds to vars
- Added a feature to write logs to file
- Added prometheus metrics with key metrics

# TODO:
- <s>Move timing (wait, latency) values to the separate config file or structure</s>
- <s>Add input file/entry validation</s>
- <s>Add an option to write logs to the file</s>
- <s>Expose prometheus metrics</s>

