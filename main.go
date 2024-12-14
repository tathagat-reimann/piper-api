// create a web server
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
)

type TextToConvert struct {
	Text string
}

func main() {
	// HTML
	http.Handle("/", http.FileServer(http.Dir("./static")))

	// API
	http.HandleFunc("/converTextToVoice", converTextToVoice)
	err := http.ListenAndServe(":3333", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}

func converTextToVoice(w http.ResponseWriter, r *http.Request) {
	//log.Println(r.Body)
	decoder := json.NewDecoder(r.Body)
	var t TextToConvert
	err := decoder.Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//log.Println("Will try to convert " + t.Text + " to voice")
	// save the text to a temp file
	in_file, err := os.CreateTemp("", "text-to-convert-*.txt")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer in_file.Close()

	_, err = in_file.WriteString(t.Text)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filePath := in_file.Name()
	log.Printf("Saved text to file: %s", filePath)

	// read ENV variables
	piper_executable := os.Getenv("PIPER_EXECUTABLE")
	piper_model := os.Getenv("PIPER_MODEL_FILE_ONNX")

	// create a temp file for the output
	out_file, err := os.CreateTemp("", "output-*.wav")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// call piper to convert the text to voice
	piper_command := "cat '" + filePath + "' | " + piper_executable +
		" --model " + piper_model +
		" -f " + out_file.Name()
	_, err = exec.Command("bash", "-c", piper_command).Output()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// send the file to the client
	w.Header().Set("Content-Disposition", "attachment;")
	file_bytes, err := os.ReadFile(out_file.Name())
	if err != nil {
		//panic(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(file_bytes)
}
