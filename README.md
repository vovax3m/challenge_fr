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
docker run -ti -e LOG_LEVEL=DEBUG -v '/path/to/sample.yaml:/go/sample.yaml:ro' challenge_fr:1 /go/sample.yaml
```
# Feature flags
```
export LOG_LEVEL=DEBUG // to enable detailed information
unset LOG_LEVEL // to show only report and errors

export LOG_FILE=/path/log.txt //to enable writing logs to the file
unset LOG_FILE // to send  logs only to Stdout/Stderr
```
# Found issues with original code:

- go.mod not exists
- main.go:87 used fmt instead log
- port ignore not implemented
- latency check not implemented
- request Body took whole endpoint not endpoint.Body
- timing not implemented
- no comments or additional messages

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

# TODO:
- <s>Move timing (wait, latency) values to the separate config file or structure</s>
- <s>Add input file/entry validation</s>
- <s>Add an option to write logs to the file</s>
- Expose prometheus metrics

