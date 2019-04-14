package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {
	client := &http.Client{}
	resp, err := client.Get("https://localhost:8012/hello")
	if err != nil {
		panic("failed to connect: " + err.Error())
	}
	content, _ := ioutil.ReadAll(resp.Body)
	s := strings.TrimSpace(string(content))

	fmt.Println(s)
}
