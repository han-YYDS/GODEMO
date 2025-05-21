package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello world")
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "pong")
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	// only post
	if r.Method != http.MethodPost {
		http.Error(w, " Only Post Method", http.StatusMethodNotAllowed)
		return
	}

	// 读取body中string
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Fail to read body", http.StatusInternalServerError)
	}
	defer r.Body.Close()

	//
	fmt.Println(string(body))
	fmt.Fprintln(w, string(body))
}

func fileHandler(w http.ResponseWriter, r *http.Request) {
	// only post
	if r.Method != http.MethodPost {
		http.Error(w, " Only Post Method", http.StatusMethodNotAllowed)
		return
	}

	// filename
	filename, err := io.ReadAll(r.Body)
	filenamestr := string(filename)
	if err != nil {
		http.Error(w, "Fail to read body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	fmt.Println("filename: ", filenamestr)
	fmt.Fprintln(w, "filename: ", filenamestr)

	// open file
	file, _ := os.Open(string(filename))
	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "Fail to read body", http.StatusInternalServerError)
		return
	}

	// return data
	w.Write(fileData)
}

func main() {
	// 注册handler
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/echo", echoHandler)
	http.HandleFunc("/file", fileHandler)

	port := ":8080"
	for {
		err := http.ListenAndServe(port, nil)
		if err != nil {
			fmt.Println("ERROR", err)
		}
	}
}
