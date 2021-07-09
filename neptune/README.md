# The API for neptune

## Reading
https://www.gwos.com/application-monitoring-with-the-prometheus-client-and-groundwork-monitor/  

https://linuxize.com/post/curl-post-request/  

https://golangdocs.com/golang-mux-router  

https://gabrieltanner.org/blog/collecting-prometheus-metrics-in-golang  

## Testing  
### Linux
#### New Temperature
`curl -X POST -F 'Temperature=68' -F 'Tank=master' http://localhost:8000/temperature`  
#### Get Temperature
`curl http://localhost:8000/temperature`  
#### Get Metrics
`curl http://localhost:8000/metrics`  

### Windows  
#### New Temperature  
`$masterTank = @{tank='master'; temperature='68'}`  
`$(Invoke-WebRequest -UseBasicParsing -Uri http://localhost:8000/temperature -Body $masterTank ).Content`  

#### Get Temperature
`$(Invoke-WebRequest -UseBasicParsing -Uri http://localhost:8000/temperature).Content`  

#### Get Metrics
`$(Invoke-WebRequest -UseBasicParsing -Uri http://localhost:8000/metrics).Content`  