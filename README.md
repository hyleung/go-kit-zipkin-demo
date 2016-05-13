#Go-Kit + Zipkin Demo

"Hello World" [Go-Kit](https://github.com/go-kit/kit) web service with [Zipkin](http://zipkin.io) tracing enabled.

## Getting Started

This project uses [GB](https://getgb.io), so it should be fairly self-contained. Assuming, you
already have `gb` installed...

`gb build all`

Start Zipkin via `docker-compose`:

`docker-compose up`

Start the service:

`bin/echoservice -port=<some port> -scribeHost=<some scribe host> -sampleRate=<something between 0 and 1>`


