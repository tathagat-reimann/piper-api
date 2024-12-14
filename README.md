This exposes an endpoint which calls Piper to convert text to voice.

Piper must be installed.
TODO: location of piper via property/yml

To build an executable for linux:
env GOOS=linux GOARCH=amd64 go build -v github.com/tathagat-reimann/piper-api
