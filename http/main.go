package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
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

func test1() {
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

const (
	MAX_UPLOAD_SIZE = 1024 * 1024 * 20 //50MB
)

func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()

	router.POST("/upload", uploadHandler)

	return router
}

func uploadHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		log.Printf("File is too big")
		return
	}
	file, headers, err := r.FormFile("file")
	if err != nil {
		log.Printf("Error when try to get file: %v", err)
		return
	}
	//获取上传文件的类型
	if headers.Header.Get("Content-Type") != "image/png" {
		log.Printf("只允许上传png图片")
		return
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("Read file error: %v", err)
		return
	}
	fn := headers.Filename
	err = ioutil.WriteFile("./video/"+fn, data, 0666)
	if err != nil {
		log.Printf("Write file error: %v", err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, "Uploaded successfully")
}

func test2() {
	r := RegisterHandlers()

	http.ListenAndServe(":8080", r)
}

func main() {
	test2()
}
