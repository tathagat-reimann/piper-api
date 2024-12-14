This exposes an endpoint which calls Piper to convert text to voice.

Piper must be installed.
TODO: location of piper via property/yml

To build an executable for linux:
env GOOS=linux GOARCH=amd64 go build -v github.com/tathagat-reimann/piper-api

You might need to chmod +x piper-api

Once you have the executable: piper-api, note the location
Run like this:
PIPER_EXECUTABLE=%PIPER_EXECUTABLE_LOCATION% PIPER_OUTPUT_DIR=%PIPER_OUTPUT_DIR% PIPER_MODEL_FILE_ONNX=%PIPER_MODEL_FILE_ONNX% ./piper-api

Example curl:
curl --output output.wav  --request POST --data '{"Text":"This text will be converted!"}' http://localhost:3333/converTextToVoice

Point your browser to localhost:3333 (or IP:3333) to use the GUI