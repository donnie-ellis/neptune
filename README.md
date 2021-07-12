# Neptune
This application listens on an API for for a temperature to be shared and then shares it out as a prometheus metric.  Additionally it publishes the most recent temperature on the API.  It contains two pieces, neptune is the server and trident is the client that publishes the temperature.

## CI/CD
Uses drone to push a new neptune container to the dmellis hosted registry.  
