package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	file, err := ioutil.ReadFile("out.gif")
	if err != nil {
		fmt.Print(err)
	}
	w.Write(file)
}

func main() {
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		fmt.Print(err)
	}
}
