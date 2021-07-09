# The API for neptune

## Environment Variables
* **NeptunePort**
    The port that Neptune will listen on.  
* **NeptuneKey**
    The key that clients will use to connect.

## Reading
I was using this application to learn a bit about go so these are some resources I used.  
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