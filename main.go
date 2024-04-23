package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

/*
	curl -F "payload=@path-to/file.txt" http://localhost:3333/upload
*/

const LISTEN_ON = ":3333"
const API_ROUTE_POST = "/upload"
const UPLOAD_DIR = "./upload"
const MAX_UPLOAD_MB = 10

func uploadFile(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(MAX_UPLOAD_MB << 20)

	file, handler, err := r.FormFile("payload")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()

	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	dst, err := os.Create(filepath.Join(UPLOAD_DIR, handler.Filename))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error creating file:", err)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error copying file:", err)
		return
	}

	fmt.Fprintf(w, "Successfully Uploaded File\n")
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		http.Error(w, "pls POST here", http.StatusFound)
	case "POST":
		uploadFile(w, r)
	}
}

func main() {

	if _, err := os.Stat(UPLOAD_DIR); os.IsNotExist(err) {
		fmt.Println(UPLOAD_DIR, "does not exist")
		err := os.MkdirAll(md.folders[dir], os.ModePerm)
		if err != nil {
			fmt.Println("Error creating local upload dir:", err)
			return
		}
	}

	fmt.Println("Listening at " + LISTEN_ON + API_ROUTE_POST)
	fmt.Println("Uploading to " + UPLOAD_DIR)
	http.HandleFunc(API_ROUTE_POST, uploadHandler)
	http.ListenAndServe(LISTEN_ON, nil)
}
