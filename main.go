package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"html/template"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path"
	"time"
)

var secretFilePath string

func main() {
	// The flag is created in a random location with a new value every time the container is restarted
	path, err := createFlag()
	if err != nil {
		panic(err)
	}

	secretFilePath = path

	http.HandleFunc("/", ping)
	http.HandleFunc("/validate", validate)
	http.ListenAndServe(":8080", nil)
}

// ping sends a ping request to the requested web address
func ping(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)

	} else if r.Method == http.MethodPost {
		requestBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		cmdString := fmt.Sprintf("ping -c 2 %s", string(requestBytes))

		cmd := exec.Command("sh", "-c", cmdString)

		outputBytes, err := cmd.CombinedOutput()
		if err != nil {
			errorMessage := fmt.Sprintf("unable to ping address:\n%v\n%s\n", err, outputBytes)
			http.Error(w, errorMessage, http.StatusBadRequest)
			return
		}

		w.Write(outputBytes)
	}
}

// validate checks the flag file's content with the request body to see if they match
func validate(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	fileBytes, err := os.ReadFile(secretFilePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if string(fileBytes) == string(bodyBytes) {
		w.Write([]byte("Success! You found the flag."))
	} else {
		w.Write([]byte("The string submitted doesn't match the content of the flag file."))
	}
}

// createFlag will select a directory in root at random and place a flag file with a randomly-generated string
func createFlag() (string, error) {
	rootDir, err := os.ReadDir("/")
	if err != nil {
		return "", err
	}

	randBytes := getRandomByteArray(64)

	// Integer of the random bytes
	randNum := binary.BigEndian.Uint64(randBytes)

	// Random number corresponding to the number of directories in root
	folderPosition := int(randNum) % (len(rootDir) - 1)

	// Placing the flag in a random directory in root
	path := path.Join("/", rootDir[folderPosition].Name(), "flag_to_capture")

	outfile, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return "", err
	}

	if _, err := io.Copy(outfile, bytes.NewReader(randBytes)); err != nil {
		return "", err
	}
	outfile.Close()

	return path, nil
}

// Adapted from https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
func getRandomByteArray(length int) []byte {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)

	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return b
}
