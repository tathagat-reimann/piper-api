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
	"strings"
	"time"
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
	file, err := os.CreateTemp("", "text-to-convert-*.txt")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	_, err = file.WriteString(t.Text)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filePath := file.Name()
	log.Printf("Saved text to file: %s", filePath)

	// call piper to convert the text to voice
	piper_executable := os.Getenv("PIPER_EXECUTABLE")
	piper_output_dir := os.Getenv("PIPER_OUTPUT_DIR")
	piper_model := os.Getenv("PIPER_MODEL_FILE_ONNX")

	randomFileName := time.Now().Format("20060102150405")

	piper_command := "cat '" + filePath + "' | " + piper_executable +
		" --model " + piper_model +
		" -f " + piper_output_dir + "/" + randomFileName + ".wav"
	out, err := exec.Command("bash", "-c", piper_command).Output()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	generated_file := strings.TrimSpace(string(out))

	log.Printf("Piper generated voice file: %s", generated_file)

	/*
		stat, _ := os.Stat(generated_file)
		log.Println(stat)
	*/

	//http.ServeFile(w, r, string(out))
	w.Header().Set("Content-Disposition", "attachment;")
	file_bytes, err := os.ReadFile(generated_file)
	if err != nil {
		//panic(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(file_bytes)
}
