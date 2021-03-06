package main

import (
	"log"
	"net/http"
)

func HelloServer(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This is an example server.\n"))
}

func main() {
	http.HandleFunc("/hello", HelloServer)
	err := http.ListenAndServeTLS(":8012", "./conf/server.pem", "./conf/server.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
