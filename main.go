// Upload a file to the server (typically on /tmp)
// Author: mfrw -- 17-07-2017
// The good old IIIT-Days

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var (
	port = flag.String("port", ":8080", "host:port to listen on")
)

func handler(res http.ResponseWriter, req *http.Request) {
	// response POST
	if req.Method == "POST" {
		src, hdr, err := req.FormFile("my-file")
		if err != nil {
			http.Error(res, err.Error(), 500)
			return
		}
		defer src.Close()

		dst, err := os.Create(filepath.Join(os.TempDir(), hdr.Filename))
		if err != nil {
			http.Error(res, err.Error(), 500)
			return
		}
		defer dst.Close()
		fmt.Println("File: ", hdr.Filename)

		io.Copy(dst, src)
		return
	}

	// Initial get request
	res.Header().Set("Content-Type", "text/html")
	hostname, err := os.Hostname()

	if err != nil {
		hostname = "localhost"
	}

	// sloppy HTML :(
	form := `
	<html>
		<head>
			<title> Upload file</title>
		</head>
		<body>
			<h1> Upload file to ` + os.TempDir() + ` on ` + hostname + `</h1>
			<form method="POST" enctype="multipart/form-data"><input type="file" name="my-file">
			<input type="submit">
			</form>
		</body>
	</html>`
	io.WriteString(res, form)
}

func main() {
	flag.Parse()
	log.Printf("[+] Server started on port: %v\n", *port)
	http.HandleFunc("/", handler)
	http.ListenAndServe(*port, nil)
}
