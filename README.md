This exposes an endpoint and offers a simple GUI, which calls Piper to convert text to voice.

## Prerequisites

Piper must be installed.

## Building the Executable

To build an executable for Linux:
```sh
env GOOS=linux GOARCH=amd64 go build -v github.com/tathagat-reimann/piper-api
```

You might need to make the executable file runnable:
```sh
chmod +x piper-api
```

## Running the Executable

Once you have the executable (`piper-api`), note its location. Also, note the location of the model file. Run the executable like this:
```sh
PIPER_EXECUTABLE=%PIPER_EXECUTABLE_LOCATION% PIPER_MODEL_FILE_ONNX=%PIPER_MODEL_FILE_ONNX% ./piper-api
```

## Example Usage

Example `curl` command:
```sh
curl --output output.wav --request POST --data '{"Text":"This text will be converted!"}' http://localhost:3333/converTextToVoice
```

## GUI Access

Point your browser to `localhost:3333` (or `IP:3333`) to use the GUI.
