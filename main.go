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
)

type TextToConvert struct {
	Text string
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./js/public")))
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
	log.Println(r.Body)
	decoder := json.NewDecoder(r.Body)
	var t TextToConvert
	err := decoder.Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("Will try to convert " + t.Text + " to voice")

	piper_executable := os.Getenv("PIPER_EXECUTABLE")

	piper_command := "echo '" + t.Text + "' | " + piper_executable + " --model /home/tathagat/tmp/piper/voices/en_US-hfc_male-medium.onnx -d /home/tathagat/tmp/piper/output/"
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
	file, err := os.ReadFile(generated_file)
	if err != nil {
		//panic(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(file)
}
